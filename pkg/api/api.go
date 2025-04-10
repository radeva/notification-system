package api

import (
	"net/http"
	"notification-system/pkg/model"
	"notification-system/pkg/queue"
	"notification-system/pkg/storage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	db    *storage.Database
	queue *queue.QueueClient
}

func NewServer(db *storage.Database, q *queue.QueueClient) *Server {
	return &Server{
		db:    db,
		queue: q,
	}
}

func (s *Server) Start() {
	r := gin.Default()

	r.POST("/notifications", func(c *gin.Context) {
		var req model.Notification
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.ID = uuid.New().String()
		req.Status = model.StatusPending
		req.CreatedAt = time.Now()

		if err := s.db.SaveNotification(req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save notification"})
			return
		}

		if err := s.queue.Publish(req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue message"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"id": req.ID, "status": req.Status})
	})

	r.GET("/notifications/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		notification, err := s.db.GetNotificationByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusOK, notification)
	})

	r.Run()
}