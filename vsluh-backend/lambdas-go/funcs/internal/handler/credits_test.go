package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// Списание — операция владельца баланса: без секрета она недоступна,
// а с настроенным auth, но без токена, отвечает 401 и не трогает DynamoDB.
func TestSpendCreditRequiresAuth(t *testing.T) {
	jwtSecret = nil
	usersTable = ""

	resp, err := SpendCredit(context.Background(), events.APIGatewayProxyRequest{})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", resp.StatusCode)
	}

	jwtSecret = []byte("test-secret")
	usersTable = "users"
	t.Cleanup(func() { jwtSecret, usersTable = nil, "" })

	resp, err = SpendCredit(context.Background(), events.APIGatewayProxyRequest{})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", resp.StatusCode)
	}
}
