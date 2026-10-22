package domain

import (
	"errors"
	"regexp"
	"strings"
)

// MinPasswordLen дублирует ограничение формы на фронте (см. L1, правило 2).
const MinPasswordLen = 6

// MaxNameLen ограничивает отображаемое имя, чтобы оно помещалось в интерфейс.
const MaxNameLen = 64

// SignupCredits — подарок за регистрацию: три разбора аудиодорожки.
const SignupCredits = 3

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("password is too short")
	ErrInvalidName     = errors.New("invalid name")
)

// emailRe — намеренно нестрогая проверка формы адреса: подтверждения почты нет,
// поэтому задача правила — отсечь опечатки, а не валидировать RFC 5322.
var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@.]+(\.[^\s@.]+)+$`)

// User — учётная запись пользователя. Хеш пароля наружу не отдаётся никогда.
type User struct {
	ID           string `json:"id"         dynamodbav:"user_id"`
	Email        string `json:"email"      dynamodbav:"email"`
	Name         string `json:"name"       dynamodbav:"name"`
	PasswordHash string `json:"-"          dynamodbav:"password_hash"`
	CreatedAt    int64  `json:"created_at" dynamodbav:"created_at"`
	// Credits — остаток разборов речи. У записей, созданных до появления
	// баланса, атрибута нет, и DynamoDB вернёт ноль — это верное поведение.
	Credits int `json:"credits" dynamodbav:"credits"`
}

// NormalizeEmail приводит адрес к ключу таблицы: регистр и пробелы не различаем.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizeName убирает лишние пробелы вокруг имени.
func NormalizeName(name string) string {
	return strings.TrimSpace(name)
}

// ValidateEmail проверяет уже нормализованный адрес.
func ValidateEmail(email string) error {
	if len(email) > 254 || !emailRe.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword проверяет длину пароля. bcrypt режет всё после 72 байт,
// поэтому длинные пароли отклоняем явно, а не молча обрезаем.
func ValidatePassword(password string) error {
	if len([]rune(password)) < MinPasswordLen || len(password) > 72 {
		return ErrInvalidPassword
	}
	return nil
}

// ValidateName проверяет отображаемое имя.
func ValidateName(name string) error {
	if name == "" || len([]rune(name)) > MaxNameLen {
		return ErrInvalidName
	}
	return nil
}
