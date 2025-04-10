package main

import (
	"log"
	"notification-system/pkg/api"
	"notification-system/pkg/config"
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

	server := api.NewServer(db, q, cfg)
	server.Start()
}
