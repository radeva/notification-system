package storage

import (
	"encoding/json"
	"fmt"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Set connection pool settings from config
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *Database) SaveNotification(n model.Notification) error {
	metadata, _ := json.Marshal(n.Metadata)
	_, err := d.db.NamedExec(`
		INSERT INTO notifications (id, channel, recipient, message, metadata, status, attempts, created_at)
		VALUES (:id, :channel, :recipient, :message, :metadata, :status, :attempts, :created_at)`,
		map[string]interface{}{
			"id":         n.ID,
			"channel":    n.Channel,
			"recipient":  n.Recipient,
			"message":    n.Message,
			"metadata":   metadata,
			"status":     n.Status,
			"attempts":   n.Attempts,
			"created_at": n.CreatedAt,
		})
	return err
}

func (d *Database) UpdateNotificationStatus(n model.Notification) error {
	metadata, _ := json.Marshal(n.Metadata)
	_, err := d.db.NamedExec(`
		UPDATE notifications SET status = :status, attempts = :attempts, last_error = :last_error, last_tried = :last_tried, metadata = :metadata
		WHERE id = :id`,
		map[string]interface{}{
			"id":         n.ID,
			"status":     n.Status,
			"attempts":   n.Attempts,
			"last_error": n.LastError,
			"last_tried": n.LastTried,
			"metadata":   metadata,
		})
	return err
}

func (d *Database) GetNotificationByID(id string) (*model.Notification, error) {
	var n model.Notification
	row := d.db.QueryRowx("SELECT * FROM notifications WHERE id = $1", id)
	if err := row.StructScan(&n); err != nil {
		return nil, err
	}
	return &n, nil
}