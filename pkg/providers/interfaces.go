package providers

import (
	"context"
	"notification-system/pkg/model"
)

// NotificationProvider defines the interface for sending notifications
type NotificationProvider interface {
	Send(ctx context.Context, notification model.Notification) error
}

// SMSProvider defines the interface for SMS notifications
type SMSProvider interface {
	NotificationProvider
}

// EmailProvider defines the interface for email notifications
type EmailProvider interface {
	NotificationProvider
}

// SlackProvider defines the interface for Slack notifications
type SlackProvider interface {
	NotificationProvider
} 