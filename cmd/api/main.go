package main

import (
	"log"
	"notification-system/pkg/api"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"notification-system/pkg/providers"
	"notification-system/pkg/queue"
	"notification-system/pkg/storage"
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

	slackProvider := providers.NewSlackNotificationProvider(cfg.Slack)
	notifier.RegisterStrategy(model.ChannelSlack, slackProvider)

	emailProvider := providers.NewEmailNotificationProvider(cfg.Email)
	notifier.RegisterStrategy(model.ChannelEmail, emailProvider)

	server := api.NewServer(db, q, cfg, notifier)
	server.Start()
}
