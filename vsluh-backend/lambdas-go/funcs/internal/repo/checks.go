package repo

import (
	"context"
	"fmt"

	"common/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// statsWindow — сколько последних проверок участвует в агрегатах кабинета.
const statsWindow = 500

// Checks хранит историю разборов: ключ — почта владельца, сортировка — время.
type Checks struct {
	db    *dynamodb.Client
	table string
}

func NewChecks(db *dynamodb.Client, table string) *Checks {
	return &Checks{db: db, table: table}
}

// SortKey делает ключ сортировки: время в миллисекундах с фиксированной шириной,
// иначе лексикографический порядок разъедется на смене разрядности.
func SortKey(createdAtMs int64, id string) string {
	return fmt.Sprintf("%013d#%s", createdAtMs, id)
}

// Save записывает проверку.
func (r *Checks) Save(ctx context.Context, check domain.Check) error {
	item, err := attributevalue.MarshalMap(check)
	if err != nil {
		return fmt.Errorf("marshal check: %w", err)
	}

	if _, err := r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	}); err != nil {
		return fmt.Errorf("put check: %w", err)
	}
	return nil
}

// List возвращает проверки пользователя от новых к старым.
// cursor — sort_key последней показанной записи, пустая строка начинает сначала.
func (r *Checks) List(ctx context.Context, email string, limit int32, cursor string) ([]domain.Check, string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		KeyConditionExpression: aws.String("user_email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(limit),
	}
	if cursor != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"user_email": &types.AttributeValueMemberS{Value: email},
			"sort_key":   &types.AttributeValueMemberS{Value: cursor},
		}
	}

	out, err := r.db.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("query checks: %w", err)
	}

	var checks []domain.Check
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &checks); err != nil {
		return nil, "", fmt.Errorf("unmarshal checks: %w", err)
	}

	next := ""
	if key, ok := out.LastEvaluatedKey["sort_key"].(*types.AttributeValueMemberS); ok {
		next = key.Value
	}
	return checks, next, nil
}

// ByID достаёт одну проверку. Запрос идёт по ключу владельца, поэтому чужую
// запись получить нельзя даже зная её идентификатор.
func (r *Checks) ByID(ctx context.Context, email, id string) (domain.Check, bool, error) {
	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		KeyConditionExpression: aws.String("user_email = :email"),
		FilterExpression:       aws.String("check_id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
			":id":    &types.AttributeValueMemberS{Value: id},
		},
		Limit: aws.Int32(statsWindow),
	})
	if err != nil {
		return domain.Check{}, false, fmt.Errorf("query check: %w", err)
	}
	if len(out.Items) == 0 {
		return domain.Check{}, false, nil
	}

	var check domain.Check
	if err := attributevalue.UnmarshalMap(out.Items[0], &check); err != nil {
		return domain.Check{}, false, fmt.Errorf("unmarshal check: %w", err)
	}
	return check, true, nil
}

// Recent отдаёт последние проверки для агрегатов кабинета.
func (r *Checks) Recent(ctx context.Context, email string) ([]domain.Check, error) {
	checks, _, err := r.List(ctx, email, statsWindow, "")
	return checks, err
}
