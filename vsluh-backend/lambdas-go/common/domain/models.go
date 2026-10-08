// Package domain содержит доменные модели сервиса ВСЛУХ.
// JSON- и DynamoDB-теги — snake_case, поля Go — PascalCase (см. L1, правило 3).
package domain

// Mode — режим практики: разная длительность подготовки и ответа.
type Mode string

const (
	ModeSpontaneous Mode = "spontaneous"
	ModeDeep        Mode = "deep"
)

// Durations возвращает время подготовки и время речи в секундах.
func (m Mode) Durations() (prepSec, speakSec int) {
	switch m {
	case ModeDeep:
		return 120, 180
	default:
		return 30, 60
	}
}

// Valid сообщает, поддерживается ли режим.
func (m Mode) Valid() bool {
	return m == ModeSpontaneous || m == ModeDeep
}

// Topic — тематическая категория вопросов.
type Topic struct {
	ID    string `json:"id"    dynamodbav:"topic_id"`
	Label string `json:"label" dynamodbav:"label"`
}

// Question — карточка вопроса для практики.
type Question struct {
	ID      string `json:"id"       dynamodbav:"question_id"`
	TopicID string `json:"topic_id" dynamodbav:"topic_id"`
	Text    string `json:"text"     dynamodbav:"text"`
	Mood    string `json:"mood"     dynamodbav:"mood,omitempty"`
}

// Session — сеанс практики пользователя.
type Session struct {
	ID         string   `json:"id"          dynamodbav:"session_id"`
	UserID     string   `json:"user_id"     dynamodbav:"user_id"`
	Mode       Mode     `json:"mode"        dynamodbav:"mode"`
	TopicIDs   []string `json:"topic_ids"   dynamodbav:"topic_ids"`
	WithMood   bool     `json:"with_mood"   dynamodbav:"with_mood"`
	QuestionID string   `json:"question_id" dynamodbav:"question_id"`
	CreatedAt  int64    `json:"created_at"  dynamodbav:"created_at"`
}
