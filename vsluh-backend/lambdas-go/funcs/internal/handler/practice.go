package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"common/domain"

	"github.com/aws/aws-lambda-go/events"
)

// topics — справочник тем. Пока статичный; при появлении редактуры переедет в DynamoDB.
var topics = []domain.Topic{
	{ID: "personal", Label: "Личное"},
	{ID: "people", Label: "Люди"},
	{ID: "work", Label: "Работа и деньги"},
	{ID: "life", Label: "Жизнь"},
	{ID: "opinions", Label: "Мнения"},
	{ID: "stories", Label: "Истории"},
	{ID: "fantasy", Label: "Фантазия"},
	{ID: "unexpected", Label: "Неожиданное"},
}

func topicExists(id string) bool {
	for _, t := range topics {
		if t.ID == id {
			return true
		}
	}
	return false
}

// ListTopics возвращает доступные темы практики.
func ListTopics(_ context.Context) (events.APIGatewayProxyResponse, error) {
	return respondJSON(http.StatusOK, map[string]any{"topics": topics})
}

type createSessionRequest struct {
	Mode     domain.Mode `json:"mode"`
	TopicIDs []string    `json:"topic_ids"`
	WithMood bool        `json:"with_mood"`
}

// CreateSession валидирует параметры практики. Те же правила продублированы на фронте.
func CreateSession(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body createSessionRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return respondError(http.StatusBadRequest, "invalid json body")
	}

	if !body.Mode.Valid() {
		return respondError(http.StatusUnprocessableEntity, "unsupported mode")
	}
	if len(body.TopicIDs) == 0 {
		return respondError(http.StatusUnprocessableEntity, "at least one topic is required")
	}
	for _, id := range body.TopicIDs {
		if !topicExists(id) {
			return respondError(http.StatusUnprocessableEntity, "unknown topic: "+id)
		}
	}

	prepSec, speakSec := body.Mode.Durations()
	return respondJSON(http.StatusCreated, map[string]any{
		"mode":        body.Mode,
		"topic_ids":   body.TopicIDs,
		"with_mood":   body.WithMood,
		"prepare_sec": prepSec,
		"speak_sec":   speakSec,
	})
}
