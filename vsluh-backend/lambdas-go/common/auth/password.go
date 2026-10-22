// Package auth содержит общие примитивы аутентификации: хеширование паролей
// и подпись токенов сессии.
package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ErrWrongPassword возвращается, когда пароль не совпал с хешем.
var ErrWrongPassword = errors.New("wrong password")

// HashPassword считает bcrypt-хеш пароля. Стоимость по умолчанию (10) —
// компромисс между защитой и холодным стартом лямбды на 256 МБ.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword сверяет пароль с хешем.
func VerifyPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrWrongPassword
	}
	return nil
}
