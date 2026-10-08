package router

import (
	"context"
	"net/http"
	"strings"

	"funcs/internal/handler"

	"github.com/aws/aws-lambda-go/events"
)

// JSONHeaders — общие заголовки ответа. CORS настроен на уровне API Gateway,
// здесь дублируем только то, что нужно самому ответу.
func JSONHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json; charset=utf-8"}
}

// Route разбирает метод и путь запроса и вызывает соответствующий хендлер.
func Route(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := strings.TrimSuffix(strings.TrimPrefix(req.Path, "/api"), "/")
	if path == "" {
		path = "/"
	}

	switch {
	case req.HTTPMethod == http.MethodOptions:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil

	case req.HTTPMethod == http.MethodGet && path == "/health":
		return handler.Health(ctx)

	case req.HTTPMethod == http.MethodGet && path == "/v1/topics":
		return handler.ListTopics(ctx)

	case req.HTTPMethod == http.MethodPost && path == "/v1/sessions":
		return handler.CreateSession(ctx, req)

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    JSONHeaders(),
			Body:       `{"error":"not found"}`,
		}, nil
	}
}
