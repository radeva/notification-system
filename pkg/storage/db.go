package storage

import (
	"context"
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

func (d *Database) SaveNotification(ctx context.Context, n model.Notification) error {
	metadata, _ := json.Marshal(n.Metadata)
	_, err := d.db.NamedExecContext(ctx, `
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

func (d *Database) UpdateNotificationStatus(ctx context.Context, n model.Notification) error {
	metadata, _ := json.Marshal(n.Metadata)
	_, err := d.db.NamedExecContext(ctx, `
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

func (d *Database) GetNotificationByID(ctx context.Context, id string) (*model.Notification, error) {
	var n model.Notification
	var metadataJSON []byte

	err := d.db.QueryRowxContext(ctx, `
		SELECT id::text, channel, recipient, message, metadata, status, attempts, last_error, last_tried, created_at
		FROM notifications 
		WHERE id::text = $1`, id).Scan(
		&n.ID,
		&n.Channel,
		&n.Recipient,
		&n.Message,
		&metadataJSON,
		&n.Status,
		&n.Attempts,
		&n.LastError,
		&n.LastTried,
		&n.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	// Unmarshal metadata if it exists
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &n.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &n, nil
}