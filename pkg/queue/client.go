package queue

import (
	"encoding/json"
	"fmt"
	"notification-system/pkg/config"
	"notification-system/pkg/model"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "notifications"
	dlqExchangeName = "notifications.dlq"
)

type QueueClient struct {
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

func NewQueueClient(cfg config.RabbitMQConfig) (*QueueClient, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare the main exchange
	err = ch.ExchangeDeclare(
		exchangeName, // exchange name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare main exchange: %w", err)
	}

	// Declare the DLQ exchange
	err = ch.ExchangeDeclare(
		dlqExchangeName, // exchange name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare DLQ exchange: %w", err)
	}

	// Declare all channel-specific queues
	for channel, queueName := range cfg.ChannelQueues {
		// Declare the main queue
		_, err = ch.QueueDeclare(
			queueName,
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			amqp.Table{
				"x-dead-letter-exchange": dlqExchangeName,
				"x-dead-letter-routing-key": queueName + ".dlq",
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue for channel %s: %w", channel, err)
		}

		// Declare the DLQ queue
		dlqQueueName := queueName + ".dlq"
		_, err = ch.QueueDeclare(
			dlqQueueName,
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare DLQ queue for channel %s: %w", channel, err)
		}

		// Bind main queue to main exchange
		err = ch.QueueBind(
			queueName,      // queue name
			string(channel), // routing key
			exchangeName,   // exchange
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to bind queue for channel %s: %w", channel, err)
		}

		// Bind DLQ queue to DLQ exchange
		err = ch.QueueBind(
			dlqQueueName,           // queue name
			queueName + ".dlq",     // routing key
			dlqExchangeName,        // exchange
			false,                  // no-wait
			nil,                    // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to bind DLQ queue for channel %s: %w", channel, err)
		}

		fmt.Printf("Declared and bound queues for channel %s: %s (main) and %s (DLQ)\n", channel, queueName, dlqQueueName)
	}

	return &QueueClient{
		channel: ch,
		config:  cfg,
	}, nil
}

func (q *QueueClient) Publish(msg model.Notification) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Verify channel is configured
	if _, ok := q.config.ChannelQueues[msg.Channel]; !ok {
		return fmt.Errorf("no queue configured for channel: %s", msg.Channel)
	}

	return q.channel.Publish(
		exchangeName, // exchange
		string(msg.Channel), // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
}

func (q *QueueClient) Consume(channel model.NotificationChannel) (<-chan amqp.Delivery, error) {
	queueName, ok := q.config.ChannelQueues[channel]
	if !ok {
		return nil, fmt.Errorf("no queue configured for channel: %s", channel)
	}

	return q.channel.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (q *QueueClient) Close() error {
	if q.channel != nil {
		return q.channel.Close()
	}
	return nil
}
