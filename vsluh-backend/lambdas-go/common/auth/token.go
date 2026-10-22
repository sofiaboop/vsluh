package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrMalformedToken = errors.New("malformed token")
	ErrBadSignature   = errors.New("bad token signature")
	ErrTokenExpired   = errors.New("token expired")
)

// Claims — полезная нагрузка токена сессии. Имя кладём внутрь, чтобы шапка
// рисовалась без дополнительного запроса к API.
type Claims struct {
	Sub       string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// jwtHeader захардкожен: единственный поддерживаемый алгоритм — HS256.
var jwtHeader = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

// IssueToken подписывает JWT HS256 со сроком жизни ttl.
func IssueToken(secret []byte, claims Claims, now time.Time, ttl time.Duration) (string, error) {
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(ttl).Unix()

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signing := jwtHeader + "." + base64.RawURLEncoding.EncodeToString(payload)
	return signing + "." + sign(secret, signing), nil
}

// ParseToken проверяет подпись и срок жизни токена.
func ParseToken(secret []byte, token string, now time.Time) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != jwtHeader {
		return Claims{}, ErrMalformedToken
	}

	signing := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(sign(secret, signing)), []byte(parts[2])) {
		return Claims{}, ErrBadSignature
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrMalformedToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, ErrMalformedToken
	}
	if claims.Sub == "" {
		return Claims{}, ErrMalformedToken
	}
	if claims.ExpiresAt <= now.Unix() {
		return Claims{}, ErrTokenExpired
	}
	return claims, nil
}

// BearerToken достаёт токен из значения заголовка Authorization.
func BearerToken(header string) string {
	const prefix = "bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}

func sign(secret []byte, signing string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signing))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
