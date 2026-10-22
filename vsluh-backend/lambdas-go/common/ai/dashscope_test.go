package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// liveASRResponse повторяет настоящий ответ qwen-audio-3.0-asr-flash:
// результат вложен в output.output, разбор лежит в объекте sentence.
const liveASRResponse = `{"output":{"output":{"sentence":{"begin_time":80,"end_time":4320,
"text":"ну, как бы, сегодня я хочу рассказать.","words":[
{"begin_time":80,"end_time":200,"text":"ну","punctuation":","},
{"begin_time":240,"end_time":400,"text":" как","punctuation":""},
{"begin_time":2400,"end_time":4320,"text":" рассказать","punctuation":"."}]},
"text":"ну, как бы, сегодня я хочу рассказать.","usage":{"duration":4}}}}`

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(Config{APIKey: "k", BaseURL: srv.URL, ASRModel: "asr", LLMModel: "llm"})
}

func TestTranscribeReadsNestedOutput(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		params, _ := body["parameters"].(map[string]any)
		if params["format"] != "webm" {
			t.Errorf("format не передан: %#v", params)
		}
		_, _ = w.Write([]byte(liveASRResponse))
	})

	tr, err := c.Transcribe(context.Background(), []byte("audio"), "webm", "ru")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if tr.Text != "ну, как бы, сегодня я хочу рассказать." {
		t.Fatalf("текст = %q", tr.Text)
	}
	if len(tr.Words) != 3 {
		t.Fatalf("слов = %d, ожидалось 3", len(tr.Words))
	}
	if tr.Words[0].Begin != 0.08 || tr.Words[2].End != 4.32 {
		t.Fatalf("тайминги не переведены в секунды: %#v", tr.Words)
	}
	if tr.Words[1].Text != "как" {
		t.Fatalf("слово не очищено от пробела: %q", tr.Words[1].Text)
	}
}

func TestTranscribeMapsSilenceToErrNoSpeech(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":"CLIENT_ERROR","message":"ASR_RESPONSE_HAVE_NO_WORDS"}`))
	})

	if _, err := c.Transcribe(context.Background(), []byte("audio"), "wav", "ru"); err != ErrNoSpeech {
		t.Fatalf("ошибка = %v, ожидалась ErrNoSpeech", err)
	}
}

func TestNotConfiguredClient(t *testing.T) {
	c := New(Config{})
	if c.ASRReady() || c.LLMReady() {
		t.Fatal("пустая конфигурация не должна считаться готовой")
	}
	if _, err := c.Transcribe(context.Background(), nil, "wav", "ru"); err != ErrNotConfigured {
		t.Fatalf("ошибка = %v, ожидалась ErrNotConfigured", err)
	}
	if _, err := c.Advise(context.Background(), "prompt"); err != ErrNotConfigured {
		t.Fatalf("ошибка = %v, ожидалась ErrNotConfigured", err)
	}
}

func TestAdviseReturnsContent(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"  Хорошо держишь мысль.  "}}]}`))
	})

	got, err := c.Advise(context.Background(), "prompt")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got != "Хорошо держишь мысль." {
		t.Fatalf("совет = %q", got)
	}
}
