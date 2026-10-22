package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"funcs/internal/repo"

	"github.com/aws/aws-lambda-go/events"
)

// SpendCredit списывает один кредит за разбор аудиодорожки.
// Разбор доступен только своему владельцу, поэтому гостям здесь делать нечего.
func SpendCredit(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !authReady() {
		return respondError(http.StatusServiceUnavailable, "auth is not configured")
	}

	claims, ok := claimsFromRequest(req)
	if !ok {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}

	credits, err := users.SpendCredit(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repo.ErrNoCredits) {
			// 402: кредиты кончились. Фронт на этот ответ показывает предложение
			// пополнить баланс и не запускает разбор.
			return respondError(http.StatusPaymentRequired, "not enough credits")
		}
		log.Printf("spend credit: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	return respondJSON(http.StatusOK, map[string]any{"credits": credits})
}
