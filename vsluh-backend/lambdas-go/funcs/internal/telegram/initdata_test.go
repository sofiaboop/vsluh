package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
	"time"
)

const testBotToken = "123456:TEST-BOT-TOKEN"

func signedInitData(authDate time.Time) string {
	params := url.Values{}
	params.Set("auth_date", fmt.Sprintf("%d", authDate.Unix()))
	params.Set("query_id", "AAE")
	params.Set("user", `{"id":42,"first_name":"Роман","username":"roman"}`)

	secret := hmacSHA256([]byte("WebAppData"), []byte(testBotToken))
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte("auth_date=" + params.Get("auth_date") + "\nquery_id=" + params.Get("query_id") + "\nuser=" + params.Get("user")))
	params.Set("hash", hex.EncodeToString(mac.Sum(nil)))

	return params.Encode()
}

func TestParseAcceptsValidSignature(t *testing.T) {
	now := time.Now()
	data, err := Parse(signedInitData(now), testBotToken, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.User.ID != 42 {
		t.Errorf("user id = %d, want 42", data.User.ID)
	}
}

func TestParseRejectsTamperedData(t *testing.T) {
	now := time.Now()
	raw := signedInitData(now) + "&extra=1"
	if _, err := Parse(raw, testBotToken, now); err == nil {
		t.Fatal("expected signature mismatch, got nil")
	}
}

func TestParseRejectsStaleData(t *testing.T) {
	authDate := time.Now().Add(-MaxAge - time.Hour)
	if _, err := Parse(signedInitData(authDate), testBotToken, time.Now()); err != ErrExpired {
		t.Fatalf("err = %v, want %v", err, ErrExpired)
	}
}

func TestParseRequiresBotToken(t *testing.T) {
	if _, err := Parse("auth_date=1&hash=aa", "", time.Now()); err != ErrNoBotToken {
		t.Fatalf("err = %v, want %v", err, ErrNoBotToken)
	}
}
