package worker

import (
	"encoding/json"
	"fmt"
	"math"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"notification-system/pkg/queue"
	"notification-system/pkg/services"
	"notification-system/pkg/storage"
	"time"
)

type Worker struct {
	db          *storage.Database
	queue       *queue.QueueClient
	smsService  *services.SMSService
	config      config.RetryConfig
}

func NewWorker(db *storage.Database, queue *queue.QueueClient, smsService *services.SMSService, config config.RetryConfig) *Worker {
	return &Worker{
		db:         db,
		queue:      queue,
		smsService: smsService,
		config:     config,
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

		err := w.process(notification)
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
			msg.Nack(false, true) // Requeue the message after all retries failed
			continue
		}

		// Acknowledge the message after successful processing
		if err := msg.Ack(false); err != nil {
			fmt.Printf("Failed to acknowledge message for channel %s: %v\n", channel, err)
		}
	}
}

func (w *Worker) process(notification model.Notification) error {
	switch notification.Channel {
	case model.ChannelEmail:
		fmt.Print("processing email")
		return nil
	case model.ChannelSMS:
		// return w.smsService.SendNotification(notification)
		fmt.Print("processing sms")
		return nil
	case model.ChannelSlack:
		fmt.Print("processing slack")
		return nil
	default:
		return fmt.Errorf("unknown channel: %s", notification.Channel)
	}
}
