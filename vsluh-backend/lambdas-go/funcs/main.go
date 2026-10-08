// Основная REST-лямбда ВСЛУХ. Проксирует все запросы API Gateway во внутренний роутер.
package main

import (
	"context"
	"log"
	"net/http"

	"funcs/internal/handler"
	"funcs/internal/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	if err := handler.Init(context.Background()); err != nil {
		log.Fatalf("handler init: %v", err)
	}
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := router.Route(ctx, req)
	if err != nil {
		log.Printf("route %s %s: %v", req.HTTPMethod, req.Path, err)
		if resp.StatusCode == 0 {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    router.JSONHeaders(),
				Body:       `{"error":"internal server error"}`,
			}, nil
		}
	}
	return resp, nil
}
