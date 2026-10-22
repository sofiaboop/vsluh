package auth

import (
	"errors"
	"testing"
	"time"
)

var secret = []byte("test-secret")

func TestIssueAndParseToken(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	token, err := IssueToken(secret, Claims{Sub: "u1", Email: "user@example.com", Name: "Роман"}, now, time.Hour)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	claims, err := ParseToken(secret, token, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.Sub != "u1" || claims.Email != "user@example.com" || claims.Name != "Роман" {
		t.Fatalf("claims = %+v", claims)
	}
	if claims.IssuedAt != now.Unix() || claims.ExpiresAt != now.Add(time.Hour).Unix() {
		t.Fatalf("время в claims неверное: %+v", claims)
	}
}

func TestParseTokenRejectsTamperedPayload(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	token, _ := IssueToken(secret, Claims{Sub: "u1"}, now, time.Hour)

	if _, err := ParseToken([]byte("other-secret"), token, now); !errors.Is(err, ErrBadSignature) {
		t.Fatalf("чужой секрет: %v, want ErrBadSignature", err)
	}
	if _, err := ParseToken(secret, token+"x", now); !errors.Is(err, ErrBadSignature) {
		t.Fatalf("испорченная подпись: %v, want ErrBadSignature", err)
	}
	if _, err := ParseToken(secret, "not-a-token", now); !errors.Is(err, ErrMalformedToken) {
		t.Fatalf("мусор вместо токена: %v, want ErrMalformedToken", err)
	}
}

func TestParseTokenExpired(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	token, _ := IssueToken(secret, Claims{Sub: "u1"}, now, time.Hour)

	if _, err := ParseToken(secret, token, now.Add(2*time.Hour)); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("протухший токен: %v, want ErrTokenExpired", err)
	}
}

func TestBearerToken(t *testing.T) {
	cases := map[string]string{
		"Bearer abc":  "abc",
		"bearer abc":  "abc",
		"BEARER  abc": "abc",
		"abc":         "",
		"":            "",
		"Basic abc":   "",
	}
	for header, want := range cases {
		if got := BearerToken(header); got != want {
			t.Errorf("BearerToken(%q) = %q, want %q", header, got, want)
		}
	}
}
