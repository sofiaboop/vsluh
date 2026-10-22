// Package storage хранит записи ответов в S3. Бакет приватный: наружу отдаётся
// только подписанная ссылка с коротким сроком жизни.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// linkTTL — срок жизни подписанной ссылки на запись. Плееру хватает,
// а утёкшая ссылка быстро перестаёт работать.
const linkTTL = 30 * time.Minute

// Audio — доступ к бакету с записями ответов.
type Audio struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

func NewAudio(client *s3.Client, bucket string) *Audio {
	a := &Audio{client: client, bucket: bucket}
	if client != nil {
		a.presign = s3.NewPresignClient(client)
	}
	return a
}

// Ready сообщает, настроено ли хранилище.
func (a *Audio) Ready() bool {
	return a != nil && a.bucket != "" && a.client != nil
}

// Key собирает путь объекта. Почта в ключ не попадает: записи раскладываем
// по идентификатору проверки, а связь с пользователем держит DynamoDB.
func Key(userID, checkID, format string) string {
	return fmt.Sprintf("checks/%s/%s.%s", userID, checkID, format)
}

// Put кладёт запись в бакет.
func (a *Audio) Put(ctx context.Context, key, contentType string, body []byte) error {
	if !a.Ready() {
		return fmt.Errorf("audio bucket is not configured")
	}

	_, err := a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("put audio: %w", err)
	}
	return nil
}

// Link выдаёт временную ссылку на запись.
func (a *Audio) Link(ctx context.Context, key string) (string, error) {
	if !a.Ready() || key == "" {
		return "", nil
	}

	out, err := a.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(linkTTL))
	if err != nil {
		return "", fmt.Errorf("presign audio: %w", err)
	}
	return out.URL, nil
}
