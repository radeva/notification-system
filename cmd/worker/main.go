package main

import (
	"log"
	"notification-system/pkg/config"
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

	// Initialize providers
	smsProvider := providers.NewTwilioSMSProvider(cfg.Twilio)
	// TODO: Initialize email and Slack providers when implemented
	
	w := worker.NewWorker(
		db,
		q,
		smsProvider,
		nil, // Email provider not implemented yet
		nil, // Slack provider not implemented yet
		cfg.Retry,
		cfg.RabbitMQ.DLQPrefix,
	)
	
	w.Start()
}