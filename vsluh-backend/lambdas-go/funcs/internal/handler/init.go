package handler

import (
	"context"
	"log"
	"os"
	"time"

	"funcs/internal/repo"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// tokenTTL — срок жизни токена сессии. Подтверждения почты нет и продления тоже:
// по истечении пользователь входит заново.
const tokenTTL = 30 * 24 * time.Hour

var (
	botToken   string
	users      *repo.Users
	jwtSecret  []byte
	usersTable string
)

// Init выполняется на холодном старте: читает конфигурацию и поднимает зависимости.
// Ошибка здесь роняет лямбду до начала обслуживания запросов — это осознанно.
func Init(ctx context.Context) error {
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	jwtSecret = []byte(os.Getenv("AUTH_JWT_SECRET"))
	usersTable = os.Getenv("DYNAMODB_USERS_TABLE_NAME")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	users = repo.NewUsers(dynamodb.NewFromConfig(cfg), usersTable)

	// Без секрета подписывать токены нечем. Роняем не всю лямбду, а только вход:
	// health и список тем должны работать и на недонастроенном окружении.
	if !authReady() {
		log.Print("auth отключён: нужны AUTH_JWT_SECRET и DYNAMODB_USERS_TABLE_NAME")
	}
	return nil
}

// authReady сообщает, может ли лямбда обслуживать вход и регистрацию.
func authReady() bool {
	return len(jwtSecret) > 0 && usersTable != ""
}
