package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	"common/auth"

	"github.com/aws/aws-lambda-go/events"
)

func TestHeaderValueIsCaseInsensitive(t *testing.T) {
	req := events.APIGatewayProxyRequest{Headers: map[string]string{"authorization": "Bearer abc"}}
	if got := headerValue(req, "Authorization"); got != "Bearer abc" {
		t.Fatalf("headerValue = %q", got)
	}
	if got := headerValue(req, "X-Missing"); got != "" {
		t.Fatalf("headerValue(missing) = %q", got)
	}
}

func TestClaimsFromRequest(t *testing.T) {
	jwtSecret = []byte("test-secret")
	t.Cleanup(func() { jwtSecret = nil })

	token, err := auth.IssueToken(jwtSecret, auth.Claims{Sub: "u1", Email: "user@example.com"}, time.Now(), time.Hour)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	req := events.APIGatewayProxyRequest{Headers: map[string]string{"Authorization": "Bearer " + token}}
	claims, ok := claimsFromRequest(req)
	if !ok || claims.Sub != "u1" {
		t.Fatalf("claimsFromRequest = %+v, ok=%v", claims, ok)
	}

	for _, header := range []string{"", "Bearer", "Bearer broken.token.value"} {
		bad := events.APIGatewayProxyRequest{Headers: map[string]string{"Authorization": header}}
		if _, ok := claimsFromRequest(bad); ok {
			t.Errorf("заголовок %q не должен давать claims", header)
		}
	}
}

// Пока секрет не задан, вход обязан отвечать 503, а не падать на nil-зависимостях.
func TestAuthDisabledWithoutSecret(t *testing.T) {
	jwtSecret = nil
	usersTable = ""

	for _, call := range []func() (events.APIGatewayProxyResponse, error){
		func() (events.APIGatewayProxyResponse, error) {
			return Register(context.Background(), events.APIGatewayProxyRequest{})
		},
		func() (events.APIGatewayProxyResponse, error) {
			return Login(context.Background(), events.APIGatewayProxyRequest{})
		},
		func() (events.APIGatewayProxyResponse, error) {
			return Me(context.Background(), events.APIGatewayProxyRequest{})
		},
	} {
		resp, err := call()
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want 503", resp.StatusCode)
		}
	}
}
