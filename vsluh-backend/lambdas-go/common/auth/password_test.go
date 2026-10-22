package auth

import (
	"errors"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "secret123" {
		t.Fatal("пароль не должен храниться в открытом виде")
	}
	if err := VerifyPassword(hash, "secret123"); err != nil {
		t.Fatalf("VerifyPassword: %v", err)
	}
	if err := VerifyPassword(hash, "secret124"); !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("VerifyPassword(wrong) = %v, want ErrWrongPassword", err)
	}
}

func TestHashPasswordIsSalted(t *testing.T) {
	first, _ := HashPassword("secret123")
	second, _ := HashPassword("secret123")
	if first == second {
		t.Fatal("хеши одинакового пароля должны различаться солью")
	}
}
