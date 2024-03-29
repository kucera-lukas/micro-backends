package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/rabbitmq"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/model"
)

const (
	providerName = "POSTGRES"
)

type messageRepository struct {
	pgxPool        *pgxpool.Pool
	rabbitmqClient *rabbitmq.Client
}

// NewMessageRepository returns implementation of the controller.Message interface.
func NewMessageRepository( // nolint:ireturn
	pgxPool *pgxpool.Pool,
	rabbitmqClient *rabbitmq.Client,
) controller.Message {
	return &messageRepository{
		pgxPool:        pgxPool,
		rabbitmqClient: rabbitmqClient,
	}
}

func (r *messageRepository) Get(
	ctx context.Context,
	messageID string,
) (*model.Message, error) {
	var message model.Message

	idInt, err := strconv.Atoi(messageID)
	if err != nil {
		return nil, fmt.Errorf(
			"get: failed to parse ID %q: %w",
			messageID,
			err,
		)
	}

	row := r.pgxPool.QueryRow(
		ctx,
		`
SELECT messages.id, messages.data, messages.created, messages.modified
FROM messages
WHERE messages.id = ($1);`,
		idInt,
	)

	if err := row.Scan(
		&message.ID,
		&message.Data,
		&message.Created,
		&message.Modified,
	); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &message, nil
}

func (r *messageRepository) Create(
	ctx context.Context,
	data string,
) (*model.Message, error) {
	var message model.Message

	row := r.pgxPool.QueryRow(
		ctx,
		`
INSERT INTO messages (data)
VALUES ($1)
RETURNING messages.id, messages.data, messages.created, messages.modified;`,
		data,
	)

	if err := row.Scan(
		&message.ID,
		&message.Data,
		&message.Created,
		&message.Modified,
	); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &message, nil
}

func (r *messageRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	row := r.pgxPool.QueryRow(ctx, "SELECT count(*) FROM messages;")
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *messageRepository) List(
	ctx context.Context,
) ([]*model.Message, error) {
	var data []*model.Message

	rows, err := r.pgxPool.Query(
		ctx,
		`
SELECT messages.id, messages.data, messages.created, messages.modified
FROM messages
LIMIT 100;`,
	)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	for rows.Next() {
		var msg model.Message

		err := rows.Scan(&msg.ID, &msg.Data, &msg.Created, &msg.Modified)
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		data = append(data, &msg)
	}

	return data, nil
}

type messageData struct {
	Data string `json:"data"`
}

func (r *messageRepository) NewMessage(
	ctx context.Context,
	delivery amqp091.Delivery,
) {
	messageData := messageData{}
	if err := json.Unmarshal(delivery.Body, &messageData); err != nil {
		log.Printf("consume: failed to umarshal message data: %v\n", err)
		nack(delivery)

		return
	}

	msg, err := r.Create(ctx, messageData.Data)
	if err != nil {
		log.Printf("consume: failed to create message: %v\n", err)
		nack(delivery)

		return
	}

	if err := r.rabbitmqClient.Publisher.Publish(
		ctx,
		fmt.Sprintf(`
{
    "message": {
        "id": %q,
        "data": %q,
        "created": %q,
        "modified": %q
    },
    "provider": %q
}`,
			strconv.Itoa(msg.ID),
			msg.Data,
			msg.Created.Format(time.RFC3339),
			msg.Modified.Format(time.RFC3339),
			providerName,
		),
		amqp091.Table{
			"provider": strings.ToLower(providerName),
			"type":     rabbitmq.CreatedMessageKey,
		},
	); err != nil {
		log.Printf(
			"consume: failed to publish message creation message: %v\n",
			err,
		)
		nack(delivery)

		return
	}

	ack(delivery)
}

func ack(delivery amqp091.Delivery) {
	if err := delivery.Ack(false); err != nil {
		log.Printf("failed to ack delivery: %v\n", err)
	}
}

func nack(delivery amqp091.Delivery) {
	if err := delivery.Nack(false, true); err != nil {
		log.Printf("failed to nack delivery: %v\n", err)
	}
}
