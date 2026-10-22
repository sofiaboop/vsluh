// Package ai — клиент Alibaba Model Studio: распознавание речи и текстовый совет.
package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ErrNotConfigured — ключ или эндпоинт не заданы, вызывать модель нечем.
var ErrNotConfigured = errors.New("model studio is not configured")

// ErrNoSpeech — модель не нашла в дорожке слов.
var ErrNoSpeech = errors.New("no speech recognized")

// Word — слово расшифровки с таймингами в секундах.
type Word struct {
	Text  string
	Begin float64
	End   float64
}

// Transcript — результат распознавания.
type Transcript struct {
	Text  string
	Words []Word
}

// Client обращается к нативному мультимодальному эндпоинту Model Studio.
// OpenAI-совместимый режим для ASR не используется: там нет таймингов слов.
type Client struct {
	apiKey   string
	baseURL  string
	asrModel string
	llmModel string
	http     *http.Client
}

// Config — настройки клиента, приходят из переменных окружения лямбды.
type Config struct {
	APIKey   string
	BaseURL  string
	ASRModel string
	LLMModel string
}

func New(cfg Config) *Client {
	return &Client{
		apiKey:   cfg.APIKey,
		baseURL:  strings.TrimSuffix(cfg.BaseURL, "/"),
		asrModel: cfg.ASRModel,
		llmModel: cfg.LLMModel,
		http:     &http.Client{Timeout: 25 * time.Second},
	}
}

// ASRReady сообщает, можно ли распознавать речь.
func (c *Client) ASRReady() bool {
	return c != nil && c.apiKey != "" && c.baseURL != "" && c.asrModel != ""
}

// LLMReady сообщает, задана ли текстовая модель. Доступ к ней может отсутствовать
// и при заданном имени — тогда Advise вернёт ошибку, а разбор обойдётся метриками.
func (c *Client) LLMReady() bool {
	return c != nil && c.apiKey != "" && c.baseURL != "" && c.llmModel != ""
}

// asrOutput повторяет форму ответа ASR. Эндпоинт отдаёт результат вложенным
// (output.output) и кладёт разбор в один объект sentence, а не в массив —
// проверено на живом ответе, доки описывают иначе.
type asrOutput struct {
	Text     string `json:"text"`
	Sentence struct {
		Text  string `json:"text"`
		Words []struct {
			Text        string  `json:"text"`
			Punctuation string  `json:"punctuation"`
			Begin       float64 `json:"begin_time"`
			End         float64 `json:"end_time"`
		} `json:"words"`
	} `json:"sentence"`
}

type asrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Output  struct {
		asrOutput
		Output asrOutput `json:"output"`
	} `json:"output"`
}

// result выбирает заполненный уровень ответа.
func (r asrResponse) result() asrOutput {
	if r.Output.Output.Text != "" || len(r.Output.Output.Sentence.Words) > 0 {
		return r.Output.Output
	}
	return r.Output.asrOutput
}

// Transcribe распознаёт дорожку. format — контейнер записи («webm», «wav», «mp3»):
// параметр обязателен, из data-URI модель его не читает.
func (c *Client) Transcribe(ctx context.Context, audio []byte, format, language string) (Transcript, error) {
	if !c.ASRReady() {
		return Transcript{}, ErrNotConfigured
	}

	payload := map[string]any{
		"model": c.asrModel,
		"input": map[string]any{
			"messages": []map[string]any{{
				"role": "user",
				"content": []map[string]any{{
					"audio": "data:audio/" + format + ";base64," + base64.StdEncoding.EncodeToString(audio),
				}},
			}},
		},
		"parameters": map[string]any{
			"format": format,
			"asr_options": map[string]any{
				"language":   language,
				"enable_lid": true,
				"enable_itn": true,
			},
		},
	}

	var out asrResponse
	if err := c.post(ctx, "/api/v1/services/aigc/multimodal-generation/generation", payload, &out); err != nil {
		// Тишина и шум — обычный исход, а не поломка: отвечаем отдельной ошибкой.
		if strings.Contains(err.Error(), "ASR_RESPONSE_HAVE_NO_WORDS") {
			return Transcript{}, ErrNoSpeech
		}
		return Transcript{}, err
	}

	res := out.result()
	tr := Transcript{Text: strings.TrimSpace(res.Text)}
	if tr.Text == "" {
		tr.Text = strings.TrimSpace(res.Sentence.Text)
	}
	for _, w := range res.Sentence.Words {
		text := strings.TrimSpace(w.Text)
		if text == "" {
			continue
		}
		// Тайминги приходят в миллисекундах.
		tr.Words = append(tr.Words, Word{Text: text, Begin: w.Begin / 1000, End: w.End / 1000})
	}
	if tr.Text == "" {
		return Transcript{}, ErrNoSpeech
	}
	return tr, nil
}

// Advise просит текстовую модель разобрать ответ. Вызов не обязателен:
// если модель недоступна, разбор остаётся на метриках и подсказках по правилам.
func (c *Client) Advise(ctx context.Context, prompt string) (string, error) {
	if !c.LLMReady() {
		return "", ErrNotConfigured
	}

	payload := map[string]any{
		"model": c.llmModel,
		"messages": []map[string]any{
			{"role": "system", "content": adviceSystemPrompt},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  400,
		"temperature": 0.6,
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := c.post(ctx, "/compatible-mode/v1/chat/completions", payload, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("model returned no choices")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

const adviceSystemPrompt = `Ты — тренер по спонтанной речи. Тебе дают расшифровку устного ответа и измерения.
Ответь по-русски, 3–5 предложений, обращайся на «ты», без списков и заголовков.
Сначала отметь сильную сторону ответа, затем назови одну главную проблему и дай конкретный приём,
как её убрать в следующий раз. Не выдумывай того, чего нет в расшифровке.`

func (c *Client) post(ctx context.Context, path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call model studio: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("model studio %d: %s", resp.StatusCode, truncate(string(raw), 300))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
