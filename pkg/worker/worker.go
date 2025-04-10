package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"notification-system/pkg/providers"
	"notification-system/pkg/queue"
	"notification-system/pkg/storage"
	"time"
)

type Worker struct {
	db          *storage.Database
	queue       *queue.QueueClient
	smsProvider providers.SMSProvider
	emailProvider providers.EmailProvider
	slackProvider providers.SlackProvider
	config      config.RetryConfig
	dlqPrefix   string
}

func NewWorker(
	db *storage.Database,
	q *queue.QueueClient,
	smsProvider providers.SMSProvider,
	emailProvider providers.EmailProvider,
	slackProvider providers.SlackProvider,
	config config.RetryConfig,
	dlqPrefix string,
) *Worker {
	return &Worker{
		db:           db,
		queue:        q,
		smsProvider:  smsProvider,
		emailProvider: emailProvider,
		slackProvider: slackProvider,
		config:       config,
		dlqPrefix:    dlqPrefix,
	}
}

func (w *Worker) Start() {
	// Start a goroutine for each channel type
	go w.processChannel(model.ChannelSMS)
	go w.processChannel(model.ChannelEmail)
	go w.processChannel(model.ChannelSlack)
	
	// Keep the main thread alive
	select {}
}

func (w *Worker) processWithRetry(notification model.Notification) error {
	var lastErr error
	for attempt := 0; attempt < w.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff
			delay := time.Duration(w.config.InitialDelayMs) * time.Millisecond * time.Duration(math.Pow(2, float64(attempt-1)))
			if delay > time.Duration(w.config.MaxDelayMs)*time.Millisecond {
				delay = time.Duration(w.config.MaxDelayMs) * time.Millisecond
			}
			time.Sleep(delay)
		}

		// Create a context with timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(w.config.ProcessTimeout)*time.Second)
		defer cancel()

		err := w.process(ctx, notification)
		if err == nil {
			return nil
		}
		lastErr = err
		fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
	}
	return fmt.Errorf("failed after %d attempts, last error: %v", w.config.MaxRetries, lastErr)
}

func (w *Worker) processChannel(channel model.NotificationChannel) {
	msgs, err := w.queue.Consume(channel)
	if err != nil {
		fmt.Printf("Failed to consume messages for channel %s: %v\n", channel, err)
		return
	}

	for msg := range msgs {
		var notification model.Notification
		if err := json.Unmarshal(msg.Body, &notification); err != nil {
			fmt.Printf("Failed to unmarshal message for channel %s: %v\n", channel, err)
			msg.Nack(false, true) // Requeue the message on error
			continue
		}

		if err := w.processWithRetry(notification); err != nil {
			fmt.Printf("Failed to process notification for channel %s after retries: %v\n", channel, err)
			
			// Send to DLQ after all retries are exhausted
			dlqName := w.dlqPrefix + string(channel)
			if err := w.queue.PublishToQueue(dlqName, msg.Body); err != nil {
				fmt.Printf("Failed to publish message to DLQ %s: %v\n", dlqName, err)
				msg.Nack(false, true) // Requeue if DLQ publish fails
				continue
			}
			
			// Acknowledge the original message since it's now in DLQ
			if err := msg.Ack(false); err != nil {
				fmt.Printf("Failed to acknowledge message for channel %s: %v\n", channel, err)
			}
			continue
		}

		// Acknowledge the message after successful processing
		if err := msg.Ack(false); err != nil {
			fmt.Printf("Failed to acknowledge message for channel %s: %v\n", channel, err)
		}
	}
}

func (w *Worker) process(ctx context.Context, notification model.Notification) error {
	switch notification.Channel {
	case model.ChannelEmail:
		if w.emailProvider == nil {
			return fmt.Errorf("email provider not configured")
		}
		return w.emailProvider.Send(ctx, notification)
	case model.ChannelSMS:
		if w.smsProvider == nil {
			return fmt.Errorf("SMS provider not configured")
		}
		return w.smsProvider.Send(ctx, notification)
	case model.ChannelSlack:
		if w.slackProvider == nil {
			return fmt.Errorf("Slack provider not configured")
		}
		return w.slackProvider.Send(ctx, notification)
	default:
		return fmt.Errorf("unknown channel: %s", notification.Channel)
	}
}
