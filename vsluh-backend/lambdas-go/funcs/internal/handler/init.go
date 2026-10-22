package handler

import (
	"context"
	"log"
	"os"
	"time"

	"common/ai"

	"funcs/internal/repo"
	"funcs/internal/storage"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// tokenTTL — срок жизни токена сессии. Подтверждения почты нет и продления тоже:
// по истечении пользователь входит заново.
const tokenTTL = 30 * 24 * time.Hour

var (
	botToken    string
	users       *repo.Users
	checks      *repo.Checks
	audioStore  *storage.Audio
	models      *ai.Client
	jwtSecret   []byte
	usersTable  string
	checksTable string
	audioBucket string
)

// Init выполняется на холодном старте: читает конфигурацию и поднимает зависимости.
// Ошибка здесь роняет лямбду до начала обслуживания запросов — это осознанно.
func Init(ctx context.Context) error {
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	jwtSecret = []byte(os.Getenv("AUTH_JWT_SECRET"))
	usersTable = os.Getenv("DYNAMODB_USERS_TABLE_NAME")
	checksTable = os.Getenv("DYNAMODB_CHECKS_TABLE_NAME")
	audioBucket = os.Getenv("AUDIO_BUCKET_NAME")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	db := dynamodb.NewFromConfig(cfg)
	users = repo.NewUsers(db, usersTable)
	checks = repo.NewChecks(db, checksTable)
	audioStore = storage.NewAudio(s3.NewFromConfig(cfg), audioBucket)

	models = ai.New(ai.Config{
		APIKey:   os.Getenv("DASHSCOPE_API_KEY"),
		BaseURL:  os.Getenv("DASHSCOPE_BASE_URL"),
		ASRModel: os.Getenv("DASHSCOPE_ASR_MODEL"),
		LLMModel: os.Getenv("DASHSCOPE_LLM_MODEL"),
	})

	// Без секрета подписывать токены нечем. Роняем не всю лямбду, а только вход:
	// health и список тем должны работать и на недонастроенном окружении.
	if !authReady() {
		log.Print("auth отключён: нужны AUTH_JWT_SECRET и DYNAMODB_USERS_TABLE_NAME")
	}
	if !checksReady() {
		log.Print("разбор отключён: нужны DYNAMODB_CHECKS_TABLE_NAME, AUDIO_BUCKET_NAME и ключ Model Studio")
	}
	if !models.LLMReady() {
		log.Print("совет текстовой модели отключён: не задан DASHSCOPE_LLM_MODEL")
	}
	return nil
}

// authReady сообщает, может ли лямбда обслуживать вход и регистрацию.
func authReady() bool {
	return len(jwtSecret) > 0 && usersTable != ""
}

// checksReady сообщает, можно ли принимать записи на разбор. Текстовая модель
// в этот список не входит: совет необязателен, метрики считаются и без него.
func checksReady() bool {
	return checksTable != "" && audioStore.Ready() && models.ASRReady()
}
