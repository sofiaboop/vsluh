package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func jsonHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json; charset=utf-8"}
}

func respondJSON(status int, payload any) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    jsonHeaders(),
		Body:       string(body),
	}, nil
}

func respondError(status int, message string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    jsonHeaders(),
		Body:       string(body),
	}, nil
}

// Health — проверка живости API.
func Health(_ context.Context) (events.APIGatewayProxyResponse, error) {
	return respondJSON(http.StatusOK, map[string]string{"status": "ok"})
}
