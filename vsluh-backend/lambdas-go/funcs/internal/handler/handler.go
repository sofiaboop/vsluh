package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Headers — заголовки любого ответа API. CORS прописан здесь, а не только
// в manifest.yaml: настройка Cors в SAM вешает заголовки лишь на preflight,
// а ответы AWS_PROXY-интеграции формирует сама лямбда.
func Headers() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json; charset=utf-8",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization,X-Telegram-Init-Data",
		"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
	}
}

func respondJSON(status int, payload any) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    Headers(),
		Body:       string(body),
	}, nil
}

func respondError(status int, message string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    Headers(),
		Body:       string(body),
	}, nil
}

// Health — проверка живости API.
func Health(_ context.Context) (events.APIGatewayProxyResponse, error) {
	return respondJSON(http.StatusOK, map[string]string{"status": "ok"})
}
