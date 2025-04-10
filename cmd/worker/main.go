package main

import (
	"log"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"notification-system/pkg/providers"
	"notification-system/pkg/queue"
	"notification-system/pkg/storage"
	"notification-system/pkg/worker"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := storage.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	q, err := queue.NewQueueClient(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Error initializing queue: %v", err)
	}
	defer q.Close()

	// Initialize notification strategy context
	notifier := providers.NewNotificationStrategyContext()
	
	// Register providers
	smsProvider := providers.NewTwilioSMSProvider(cfg.Twilio)
	notifier.RegisterStrategy(model.ChannelSMS, smsProvider)
	// TODO: Register email and Slack providers when implemented
	
	w := worker.NewWorker(
		db,
		q,
		notifier,
		cfg.Retry,
		cfg.RabbitMQ.DLQPrefix,
	)
	
	w.Start()
}