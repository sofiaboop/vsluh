// Package telegram проверяет подлинность initData Telegram Mini App.
package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// MaxAge — максимальный возраст initData. Защищает от переигрывания старых подписей.
const MaxAge = 24 * time.Hour

var (
	ErrEmptyInitData = errors.New("telegram: empty init data")
	ErrMissingHash   = errors.New("telegram: hash is missing")
	ErrInvalidHash   = errors.New("telegram: signature mismatch")
	ErrExpired       = errors.New("telegram: init data expired")
	ErrNoBotToken    = errors.New("telegram: bot token is not configured")
)

// User — пользователь Telegram из проверенного initData.
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// InitData — разобранные и проверенные данные Mini App.
type InitData struct {
	User     User
	AuthDate time.Time
}

// Parse проверяет подпись initData и возвращает данные пользователя.
// Без валидного botToken данные не принимаются — initDataUnsafe доверять нельзя.
func Parse(raw, botToken string, now time.Time) (*InitData, error) {
	if raw == "" {
		return nil, ErrEmptyInitData
	}
	if botToken == "" {
		return nil, ErrNoBotToken
	}

	params, err := url.ParseQuery(raw)
	if err != nil {
		return nil, fmt.Errorf("telegram: parse init data: %w", err)
	}

	hash := params.Get("hash")
	if hash == "" {
		return nil, ErrMissingHash
	}
	params.Del("hash")
	params.Del("signature")

	if !validSignature(params, hash, botToken) {
		return nil, ErrInvalidHash
	}

	var authDate time.Time
	if v := params.Get("auth_date"); v != "" {
		var unix int64
		if _, err := fmt.Sscanf(v, "%d", &unix); err == nil {
			authDate = time.Unix(unix, 0)
		}
	}
	if authDate.IsZero() || now.Sub(authDate) > MaxAge {
		return nil, ErrExpired
	}

	var user User
	if v := params.Get("user"); v != "" {
		if err := json.Unmarshal([]byte(v), &user); err != nil {
			return nil, fmt.Errorf("telegram: parse user: %w", err)
		}
	}
	if user.ID == 0 {
		return nil, errors.New("telegram: user id is missing")
	}

	return &InitData{User: user, AuthDate: authDate}, nil
}

func validSignature(params url.Values, hash, botToken string) bool {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+params.Get(k))
	}
	dataCheckString := strings.Join(pairs, "\n")

	secret := hmacSHA256([]byte("WebAppData"), []byte(botToken))
	expected := hmacSHA256(secret, []byte(dataCheckString))

	given, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return hmac.Equal(expected, given)
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
