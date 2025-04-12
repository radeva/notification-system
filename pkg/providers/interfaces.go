package providers

import (
	"context"
	"notification-system/pkg/model"
)

// Validator defines the interface for validating notifications
type Validator interface {
	Validate(notification model.Notification) error
}

// NotificationProvider defines the interface for sending notifications
type NotificationProvider interface {
	Validator
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