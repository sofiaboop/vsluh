package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"common/ai"
	"common/auth"
	"common/domain"

	"funcs/internal/storage"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-lambda-go/events"
)

// configureChecks включает разбор так, чтобы до сети дело не дошло:
// проверяем ранние отказы, а не обращение к модели.
func configureChecks(t *testing.T) {
	t.Helper()

	jwtSecret = []byte("test-secret")
	usersTable = "users"
	checksTable = "checks"
	audioStore = storage.NewAudio(s3.New(s3.Options{Region: "eu-central-1"}), "bucket")
	models = ai.New(ai.Config{APIKey: "key", BaseURL: "https://example.test", ASRModel: "asr"})

	t.Cleanup(func() {
		jwtSecret, usersTable, checksTable = nil, "", ""
		audioStore = storage.NewAudio(nil, "")
		models = ai.New(ai.Config{})
	})
}

func authHeader(t *testing.T) map[string]string {
	t.Helper()

	token, err := auth.IssueToken(jwtSecret, auth.Claims{Sub: "u1", Email: "a@b.co"}, time.Now(), tokenTTL)
	if err != nil {
		t.Fatalf("выпуск токена: %v", err)
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

func TestCreateCheckRequiresConfiguredAnalysis(t *testing.T) {
	jwtSecret = []byte("test-secret")
	usersTable = "users"
	checksTable = ""
	audioStore = storage.NewAudio(nil, "")
	models = ai.New(ai.Config{})
	t.Cleanup(func() { jwtSecret, usersTable = nil, "" })

	resp, err := CreateCheck(context.Background(), events.APIGatewayProxyRequest{})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", resp.StatusCode)
	}
}

func TestCreateCheckRequiresToken(t *testing.T) {
	configureChecks(t)

	resp, err := CreateCheck(context.Background(), events.APIGatewayProxyRequest{})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", resp.StatusCode)
	}
}

func TestCreateCheckRejectsUnsupportedFormat(t *testing.T) {
	configureChecks(t)

	body, _ := json.Marshal(createCheckRequest{
		AudioBase64: base64.StdEncoding.EncodeToString([]byte("audio")),
		Format:      "flac",
		Mode:        domain.ModeSpontaneous,
	})

	resp, err := CreateCheck(context.Background(), events.APIGatewayProxyRequest{
		Headers: authHeader(t),
		Body:    string(body),
	})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", resp.StatusCode)
	}
}

func TestCreateCheckRejectsOversizedAudio(t *testing.T) {
	configureChecks(t)

	body, _ := json.Marshal(createCheckRequest{
		AudioBase64: base64.StdEncoding.EncodeToString(make([]byte, domain.MaxAudioBytes+1)),
		Format:      "webm",
		Mode:        domain.ModeSpontaneous,
	})

	resp, err := CreateCheck(context.Background(), events.APIGatewayProxyRequest{
		Headers: authHeader(t),
		Body:    string(body),
	})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", resp.StatusCode)
	}
}

func TestHistoryEndpointsRequireToken(t *testing.T) {
	configureChecks(t)

	for name, call := range map[string]func() (events.APIGatewayProxyResponse, error){
		"list": func() (events.APIGatewayProxyResponse, error) {
			return ListChecks(context.Background(), events.APIGatewayProxyRequest{})
		},
		"get": func() (events.APIGatewayProxyResponse, error) {
			return GetCheck(context.Background(), events.APIGatewayProxyRequest{}, "id")
		},
		"stats": func() (events.APIGatewayProxyResponse, error) {
			return GetStats(context.Background(), events.APIGatewayProxyRequest{})
		},
	} {
		resp, err := call()
		if err != nil {
			t.Fatalf("%s: неожиданная ошибка: %v", name, err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("%s: status = %d, want 401", name, resp.StatusCode)
		}
	}
}

func TestStripDataURI(t *testing.T) {
	if got := stripDataURI("data:audio/webm;base64,QUJD"); got != "QUJD" {
		t.Fatalf("префикс не отрезан: %q", got)
	}
	if got := stripDataURI("QUJD"); got != "QUJD" {
		t.Fatalf("чистая строка изменилась: %q", got)
	}
}
