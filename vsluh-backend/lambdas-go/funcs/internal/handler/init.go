package handler

import (
	"context"
	"os"
)

var botToken string

// Init выполняется на холодном старте: читает конфигурацию и поднимает зависимости.
// Ошибка здесь роняет лямбду до начала обслуживания запросов — это осознанно.
func Init(_ context.Context) error {
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	return nil
}
