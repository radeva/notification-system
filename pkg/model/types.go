package model

import "time"

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

type NotificationChannel string

const (
	ChannelSMS   NotificationChannel = "sms"
	ChannelEmail NotificationChannel = "email"
	ChannelSlack NotificationChannel = "slack"
)

type Notification struct {
	ID        string            `db:"id" json:"id"`
	Channel   NotificationChannel `db:"channel" json:"channel"`
	Recipient string            `db:"recipient" json:"recipient"`
	Message   string            `db:"message" json:"message"`
	Metadata  map[string]string `db:"metadata" json:"metadata"`
	Status    NotificationStatus `db:"status" json:"status"`
	Attempts  int               `db:"attempts" json:"attempts"`
	LastError *string           `db:"last_error" json:"lastError,omitempty"`
	CreatedAt time.Time         `db:"created_at" json:"createdAt"`
	LastTried *time.Time        `db:"last_tried" json:"lastTried,omitempty"`
}
