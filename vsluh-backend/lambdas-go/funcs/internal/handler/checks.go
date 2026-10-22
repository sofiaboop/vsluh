package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"common/ai"
	"common/analysis"
	"common/auth"
	"common/domain"

	"funcs/internal/repo"
	"funcs/internal/storage"

	"github.com/aws/aws-lambda-go/events"
)

// historyPageSize — сколько проверок отдаём за один запрос истории.
const historyPageSize = 20

// audioFormats — контейнеры, которые присылает браузер. MediaRecorder в Chrome
// пишет webm, в Safari — mp4; остальное отсекаем, чтобы не платить за мусор.
var audioFormats = map[string]string{
	"webm": "audio/webm",
	"mp4":  "audio/mp4",
	"m4a":  "audio/mp4",
	"ogg":  "audio/ogg",
	"wav":  "audio/wav",
	"mp3":  "audio/mpeg",
}

type createCheckRequest struct {
	AudioBase64 string      `json:"audio_base64"`
	Format      string      `json:"format"`
	TopicID     string      `json:"topic_id"`
	TopicLabel  string      `json:"topic_label"`
	Mode        domain.Mode `json:"mode"`
	SpeakSec    float64     `json:"speak_sec"`
	Language    string      `json:"language"`
}

// CreateCheck принимает запись ответа, списывает кредит и возвращает разбор.
// Кредит списывается до обращения к модели, но возвращается, если разбор
// не состоялся: платить за нашу ошибку пользователь не должен.
func CreateCheck(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !authReady() {
		return respondError(http.StatusServiceUnavailable, "auth is not configured")
	}
	if !checksReady() {
		return respondError(http.StatusServiceUnavailable, "analysis is not configured")
	}

	claims, ok := claimsFromRequest(req)
	if !ok {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}

	var body createCheckRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return respondError(http.StatusBadRequest, "invalid json body")
	}

	contentType, supported := audioFormats[strings.ToLower(body.Format)]
	if !supported {
		return respondError(http.StatusUnprocessableEntity, "unsupported audio format")
	}
	if !body.Mode.Valid() {
		return respondError(http.StatusUnprocessableEntity, "unsupported mode")
	}

	audio, err := base64.StdEncoding.DecodeString(stripDataURI(body.AudioBase64))
	if err != nil || len(audio) == 0 {
		return respondError(http.StatusBadRequest, "invalid audio payload")
	}
	if len(audio) > domain.MaxAudioBytes {
		return respondError(http.StatusRequestEntityTooLarge, "audio is too large")
	}

	credits, err := users.SpendCredit(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repo.ErrNoCredits) {
			return respondError(http.StatusPaymentRequired, "not enough credits")
		}
		log.Printf("spend credit: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	check, err := buildCheck(ctx, claims.Sub, claims.Email, body, audio, contentType)
	if err != nil {
		// Разбор не состоялся — возвращаем кредит и сообщаем, что произошло.
		restored, refundErr := users.RefundCredit(ctx, claims.Email)
		if refundErr != nil {
			log.Printf("refund credit: %v", refundErr)
			restored = credits
		}
		if errors.Is(err, ai.ErrNoSpeech) {
			return respondJSON(http.StatusUnprocessableEntity, map[string]any{
				"error":   "no speech recognized",
				"credits": restored,
			})
		}
		log.Printf("build check: %v", err)
		return respondJSON(http.StatusBadGateway, map[string]any{
			"error":   "analysis failed",
			"credits": restored,
		})
	}

	if err := checks.Save(ctx, check); err != nil {
		log.Printf("save check: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	check.AudioURL = audioLink(ctx, check.AudioKey)
	return respondJSON(http.StatusCreated, map[string]any{"check": check, "credits": credits})
}

// buildCheck распознаёт запись, считает метрики и кладёт аудио в хранилище.
func buildCheck(ctx context.Context, userID, email string, body createCheckRequest, audio []byte, contentType string) (domain.Check, error) {
	language := body.Language
	if language == "" {
		language = "ru"
	}

	format := strings.ToLower(body.Format)
	transcript, err := models.Transcribe(ctx, audio, format, language)
	if err != nil {
		return domain.Check{}, err
	}

	words := make([]analysis.Word, 0, len(transcript.Words))
	for _, w := range transcript.Words {
		words = append(words, analysis.Word{Text: w.Text, Begin: w.Begin, End: w.End})
	}

	speakSec := body.SpeakSec
	if speakSec <= 0 {
		_, planned := body.Mode.Durations()
		speakSec = float64(planned)
	}

	metrics := analysis.Build(transcript.Text, words, speakSec)

	id, err := newCheckID()
	if err != nil {
		return domain.Check{}, fmt.Errorf("new check id: %w", err)
	}

	now := time.Now()
	check := domain.Check{
		ID:         id,
		UserEmail:  email,
		CreatedAt:  now.Unix(),
		SortKey:    repo.SortKey(now.UnixMilli(), id),
		TopicID:    body.TopicID,
		TopicLabel: body.TopicLabel,
		Mode:       body.Mode,
		Transcript: transcript.Text,
		Metrics:    metrics,
		Tips:       analysis.Tips(metrics),
		AudioKey:   storage.Key(userID, id, format),
	}

	if err := audioStore.Put(ctx, check.AudioKey, contentType, audio); err != nil {
		return domain.Check{}, err
	}

	check.Advice = advice(ctx, check)
	return check, nil
}

// advice просит текстовую модель прокомментировать ответ. Доступа к модели
// может не быть — тогда разбор остаётся на метриках и подсказках по правилам.
func advice(ctx context.Context, check domain.Check) string {
	if !models.LLMReady() {
		return ""
	}

	text, err := models.Advise(ctx, advicePrompt(check))
	if err != nil {
		log.Printf("advice skipped: %v", err)
		return ""
	}
	return text
}

func advicePrompt(check domain.Check) string {
	var b strings.Builder
	if check.TopicLabel != "" {
		b.WriteString("Тема: " + check.TopicLabel + "\n")
	}
	b.WriteString("Расшифровка: " + check.Transcript + "\n")
	b.WriteString("Измерения: ")
	b.WriteString("слов " + strconv.Itoa(check.Metrics.WordCount))
	b.WriteString(", темп " + strconv.Itoa(check.Metrics.WordsPerMin) + " слов в минуту")
	b.WriteString(", слов-паразитов " + strconv.Itoa(check.Metrics.FillerCount))
	b.WriteString(", повторов " + strconv.Itoa(check.Metrics.RepeatCount))
	b.WriteString(", длинных пауз " + strconv.Itoa(check.Metrics.LongPauses))
	if len(check.Metrics.FillerWords) > 0 {
		b.WriteString("\nЧаще всего звучало: ")
		for i, f := range check.Metrics.FillerWords {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(f.Word + " ×" + strconv.Itoa(f.Count))
		}
	}
	return b.String()
}

// ListChecks отдаёт историю проверок пользователя.
func ListChecks(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	claims, ok := authorizedClaims(req)
	if !ok {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}
	if !checksReady() {
		return respondError(http.StatusServiceUnavailable, "analysis is not configured")
	}

	items, next, err := checks.List(ctx, claims.Email, historyPageSize, req.QueryStringParameters["cursor"])
	if err != nil {
		log.Printf("list checks: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	if items == nil {
		items = []domain.Check{}
	}

	return respondJSON(http.StatusOK, map[string]any{"checks": items, "next_cursor": next})
}

// GetCheck отдаёт одну проверку вместе со ссылкой на запись.
func GetCheck(ctx context.Context, req events.APIGatewayProxyRequest, id string) (events.APIGatewayProxyResponse, error) {
	claims, ok := authorizedClaims(req)
	if !ok {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}
	if !checksReady() {
		return respondError(http.StatusServiceUnavailable, "analysis is not configured")
	}

	check, found, err := checks.ByID(ctx, claims.Email, id)
	if err != nil {
		log.Printf("get check: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	if !found {
		return respondError(http.StatusNotFound, "check not found")
	}

	check.AudioURL = audioLink(ctx, check.AudioKey)
	return respondJSON(http.StatusOK, map[string]any{"check": check})
}

// GetStats считает агрегаты личного кабинета.
func GetStats(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	claims, ok := authorizedClaims(req)
	if !ok {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}
	if !checksReady() {
		return respondError(http.StatusServiceUnavailable, "analysis is not configured")
	}

	items, err := checks.Recent(ctx, claims.Email)
	if err != nil {
		log.Printf("stats: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	return respondJSON(http.StatusOK, map[string]any{"stats": analysis.Aggregate(items, time.Now())})
}

// authorizedClaims проверяет и настройку входа, и сам токен.
func authorizedClaims(req events.APIGatewayProxyRequest) (auth.Claims, bool) {
	if !authReady() {
		return auth.Claims{}, false
	}
	return claimsFromRequest(req)
}

func audioLink(ctx context.Context, key string) string {
	url, err := audioStore.Link(ctx, key)
	if err != nil {
		log.Printf("audio link: %v", err)
		return ""
	}
	return url
}

// stripDataURI отрезает префикс data-URI, если фронт прислал его целиком.
func stripDataURI(payload string) string {
	if idx := strings.Index(payload, "base64,"); idx >= 0 {
		return payload[idx+len("base64,"):]
	}
	return payload
}

func newCheckID() (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
