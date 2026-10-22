package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	if got := NormalizeEmail("  User@Example.COM "); got != "user@example.com" {
		t.Fatalf("NormalizeEmail = %q", got)
	}
}

func TestValidateEmail(t *testing.T) {
	valid := []string{"user@example.com", "a.b+tag@mail.co.uk"}
	for _, email := range valid {
		if err := ValidateEmail(email); err != nil {
			t.Errorf("ValidateEmail(%q) = %v, want nil", email, err)
		}
	}

	invalid := []string{"", "user", "user@", "@example.com", "user@example", "us er@example.com",
		strings.Repeat("a", 250) + "@example.com"}
	for _, email := range invalid {
		if err := ValidateEmail(email); !errors.Is(err, ErrInvalidEmail) {
			t.Errorf("ValidateEmail(%q) = %v, want ErrInvalidEmail", email, err)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("secret"); err != nil {
		t.Errorf("ValidatePassword(6 chars) = %v, want nil", err)
	}
	if err := ValidatePassword("парол"); !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("ValidatePassword(5 runes) = %v, want ErrInvalidPassword", err)
	}
	if err := ValidatePassword(strings.Repeat("a", 73)); !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("ValidatePassword(73 bytes) = %v, want ErrInvalidPassword", err)
	}
}

func TestValidateName(t *testing.T) {
	if err := ValidateName("Роман"); err != nil {
		t.Errorf("ValidateName = %v, want nil", err)
	}
	if err := ValidateName(""); !errors.Is(err, ErrInvalidName) {
		t.Error("пустое имя должно отклоняться")
	}
	if err := ValidateName(strings.Repeat("я", MaxNameLen+1)); !errors.Is(err, ErrInvalidName) {
		t.Error("слишком длинное имя должно отклоняться")
	}
}
