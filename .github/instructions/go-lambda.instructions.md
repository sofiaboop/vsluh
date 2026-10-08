---
applyTo: "vsluh-backend/**/*.go"
---

# Go Lambda

- Каждая лямбда — отдельный Go-модуль в `lambdas-go/<name>/` с собственным `go.mod`.
- Общий код — модуль `common`, подключается через `replace common => ../common`.
- Все модули перечислены в `vsluh-backend/go.work`.

## Скелет

```go
var svc *service.TopicService // package-level: переживает тёплый старт

func init() {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        log.Fatalf("load aws config: %v", err) // холодный старт падает быстро
    }
    repo := repo.NewTopicRepo(dynamodb.NewFromConfig(cfg), os.Getenv("DYNAMODB_TOPICS_TABLE_NAME"))
    svc = service.NewTopicService(repo)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // без паник: любая ошибка -> лог + HTTP-код
}

func main() { lambda.Start(handler) }
```

## Правила

- `log.Fatalf` допустим только в `init()`. В хендлере — `log.Printf` + корректный статус.
- Имена таблиц читаются из env (`DYNAMODB_*_TABLE_NAME`), никогда не хардкодятся.
- Сборка: `GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap`.
- Тесты обязательны для `common/` (доменная логика). `make test-all` должен быть зелёным.
