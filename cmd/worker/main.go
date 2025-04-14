package main

import (
	"fmt"
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
	

	var smsProvider providers.SMSProvider
	var slackProvider providers.SlackProvider
	var emailProvider providers.EmailProvider

	if !cfg.UseMockProviders {
		smsProvider = providers.NewTwilioSMSProvider(cfg.Twilio)
		slackProvider = providers.NewSlackNotificationProvider(cfg.Slack)
		emailProvider = providers.NewEmailNotificationProvider(cfg.Email)
	} else {
		fmt.Println("NOTE: worker is using mock providers")
		smsProvider = providers.NewMockSMSProvider()
		slackProvider = providers.NewMockSlackProvider()
		emailProvider = providers.NewMockEmailProvider()
	}

	// Register providers
	notifier.RegisterStrategy(model.ChannelSMS, smsProvider)
	notifier.RegisterStrategy(model.ChannelSlack, slackProvider)
	notifier.RegisterStrategy(model.ChannelEmail, emailProvider)
	
	w := worker.NewWorker(
		db,
		q,
		notifier,
		cfg.Retry,
		cfg.RabbitMQ.DLQPrefix,
	)
	
	w.Start()
}