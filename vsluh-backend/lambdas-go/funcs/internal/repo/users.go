// Package repo — доступ к DynamoDB. Слой знает про таблицы, но не про HTTP.
package repo

import (
	"context"
	"errors"
	"fmt"

	"common/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ErrEmailTaken — адрес уже зарегистрирован.
var ErrEmailTaken = errors.New("email already registered")

// ErrNoCredits — на балансе нечего списывать.
var ErrNoCredits = errors.New("not enough credits")

// Users хранит учётные записи: ключ таблицы — нормализованный email.
type Users struct {
	db    *dynamodb.Client
	table string
}

func NewUsers(db *dynamodb.Client, table string) *Users {
	return &Users{db: db, table: table}
}

// Create сохраняет пользователя. Уникальность адреса держит условие на ключе,
// а не предварительная проверка: между чтением и записью возможна гонка.
func (r *Users) Create(ctx context.Context, user domain.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	})
	if err != nil {
		var conflict *types.ConditionalCheckFailedException
		if errors.As(err, &conflict) {
			return ErrEmailTaken
		}
		return fmt.Errorf("put user: %w", err)
	}
	return nil
}

// ByEmail возвращает пользователя по нормализованному адресу.
// Второе значение — найден ли он.
func (r *Users) ByEmail(ctx context.Context, email string) (domain.User, bool, error) {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key:       map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}},
	})
	if err != nil {
		return domain.User{}, false, fmt.Errorf("get user: %w", err)
	}
	if len(out.Item) == 0 {
		return domain.User{}, false, nil
	}

	var user domain.User
	if err := attributevalue.UnmarshalMap(out.Item, &user); err != nil {
		return domain.User{}, false, fmt.Errorf("unmarshal user: %w", err)
	}
	return user, true, nil
}

// SpendCredit списывает один кредит и возвращает остаток. Проверка баланса
// живёт в условии записи: два параллельных разбора не должны увести в минус.
func (r *Users) SpendCredit(ctx context.Context, email string) (int, error) {
	out, err := r.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           aws.String(r.table),
		Key:                 map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}},
		UpdateExpression:    aws.String("SET credits = credits - :one"),
		ConditionExpression: aws.String("attribute_exists(email) AND credits >= :one"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":one": &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	if err != nil {
		var conflict *types.ConditionalCheckFailedException
		if errors.As(err, &conflict) {
			return 0, ErrNoCredits
		}
		return 0, fmt.Errorf("spend credit: %w", err)
	}

	var updated struct {
		Credits int `dynamodbav:"credits"`
	}
	if err := attributevalue.UnmarshalMap(out.Attributes, &updated); err != nil {
		return 0, fmt.Errorf("unmarshal credits: %w", err)
	}
	return updated.Credits, nil
}
