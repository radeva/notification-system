package api

import (
	"context"
	"fmt"
	"net/http"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"notification-system/pkg/queue"
	"notification-system/pkg/storage"
	"notification-system/pkg/validation"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	db       *storage.Database
	queue    *queue.QueueClient
	cfg      *config.Config
	validator validation.Validator
}

func NewServer(db *storage.Database, q *queue.QueueClient, cfg *config.Config, validator validation.Validator) *Server {
	return &Server{
		db:       db,
		queue:    q,
		cfg:      cfg,
		validator: validator,
	}
}

func (s *Server) Start() {
	r := gin.Default()

	r.POST("/notifications", func(c *gin.Context) {
		var notification model.Notification
		if err := c.ShouldBindJSON(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if err := s.validator.Validate(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid notification: %v", err)})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(s.cfg.Server.RequestTimeout)*time.Second)
		defer cancel()

		notification.ID = uuid.New().String()
		notification.Status = model.StatusPending
		notification.CreatedAt = time.Now()

		if err := s.db.SaveNotification(ctx, notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save notification"})
			return
		}

		if err := s.queue.Publish(notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue message"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"id": notification.ID, "status": notification.Status})
	})

	r.GET("/notifications/:id/status", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(s.cfg.Server.RequestTimeout)*time.Second)
		defer cancel()

		id := c.Param("id")
		notification, err := s.db.GetNotificationByID(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusOK, notification)
	})

	r.Run()
}