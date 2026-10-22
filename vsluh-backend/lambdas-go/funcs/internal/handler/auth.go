package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"common/auth"
	"common/domain"

	"funcs/internal/repo"

	"github.com/aws/aws-lambda-go/events"
)

type credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// authResponse — единый ответ входа и регистрации: токен и профиль.
type authResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

// Register создаёт учётную запись по почте и паролю. Подтверждение адреса пока не шлём.
func Register(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !authReady() {
		return respondError(http.StatusServiceUnavailable, "auth is not configured")
	}

	var body credentials
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return respondError(http.StatusBadRequest, "invalid json body")
	}

	name := domain.NormalizeName(body.Name)
	email := domain.NormalizeEmail(body.Email)
	if err := domain.ValidateName(name); err != nil {
		return respondError(http.StatusUnprocessableEntity, "invalid name")
	}
	if err := domain.ValidateEmail(email); err != nil {
		return respondError(http.StatusUnprocessableEntity, "invalid email")
	}
	if err := domain.ValidatePassword(body.Password); err != nil {
		return respondError(http.StatusUnprocessableEntity, "password is too short")
	}

	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		log.Printf("hash password: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	id, err := newUserID()
	if err != nil {
		log.Printf("generate user id: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	user := domain.User{
		ID:           id,
		Email:        email,
		Name:         name,
		PasswordHash: hash,
		CreatedAt:    time.Now().Unix(),
		Credits:      domain.SignupCredits,
	}
	if err := users.Create(ctx, user); err != nil {
		if errors.Is(err, repo.ErrEmailTaken) {
			return respondError(http.StatusConflict, "email already registered")
		}
		log.Printf("create user: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}

	return respondWithToken(http.StatusCreated, user)
}

// Login проверяет пару почта+пароль и выдаёт токен.
func Login(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !authReady() {
		return respondError(http.StatusServiceUnavailable, "auth is not configured")
	}

	var body credentials
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return respondError(http.StatusBadRequest, "invalid json body")
	}

	user, found, err := users.ByEmail(ctx, domain.NormalizeEmail(body.Email))
	if err != nil {
		log.Printf("lookup user: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	// Один и тот же ответ на незнакомый адрес и на неверный пароль:
	// иначе форма превращается в проверку, кто зарегистрирован.
	if !found {
		return respondError(http.StatusUnauthorized, "invalid email or password")
	}
	if err := auth.VerifyPassword(user.PasswordHash, body.Password); err != nil {
		return respondError(http.StatusUnauthorized, "invalid email or password")
	}

	return respondWithToken(http.StatusOK, user)
}

// Me возвращает профиль владельца токена.
func Me(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !authReady() {
		return respondError(http.StatusServiceUnavailable, "auth is not configured")
	}

	claims, ok := claimsFromRequest(req)
	if !ok {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}

	user, found, err := users.ByEmail(ctx, claims.Email)
	if err != nil {
		log.Printf("lookup user: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	if !found {
		return respondError(http.StatusUnauthorized, "invalid or expired token")
	}

	return respondJSON(http.StatusOK, map[string]any{"user": user})
}

// claimsFromRequest разбирает токен из заголовка Authorization.
func claimsFromRequest(req events.APIGatewayProxyRequest) (auth.Claims, bool) {
	token := auth.BearerToken(headerValue(req, "Authorization"))
	if token == "" {
		return auth.Claims{}, false
	}
	claims, err := auth.ParseToken(jwtSecret, token, time.Now())
	if err != nil {
		return auth.Claims{}, false
	}
	return claims, true
}

// headerValue ищет заголовок без учёта регистра: API Gateway отдаёт их как есть.
func headerValue(req events.APIGatewayProxyRequest, name string) string {
	for key, value := range req.Headers {
		if strings.EqualFold(key, name) {
			return value
		}
	}
	return ""
}

func respondWithToken(status int, user domain.User) (events.APIGatewayProxyResponse, error) {
	token, err := auth.IssueToken(jwtSecret, auth.Claims{
		Sub:   user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, time.Now(), tokenTTL)
	if err != nil {
		log.Printf("issue token: %v", err)
		return respondError(http.StatusInternalServerError, "internal server error")
	}
	return respondJSON(status, authResponse{Token: token, User: user})
}

func newUserID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
