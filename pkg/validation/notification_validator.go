package validation

import (
	"fmt"
	"notification-system/pkg/model"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`) // E.164 format
)

// Validator defines the interface for notification validation
type Validator interface {
	Validate(notification *model.Notification) error
}

// NotificationValidator implements the Validator interface
type NotificationValidator struct{}

// NewNotificationValidator creates a new notification validator
func NewNotificationValidator() Validator {
	return &NotificationValidator{}
}

// Validate performs validation on a notification
func (v *NotificationValidator) Validate(notification *model.Notification) error {
	if notification.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	if notification.Recipient == "" {
		return fmt.Errorf("recipient cannot be empty")
	}

	switch notification.Channel {
	case model.ChannelEmail:
		if !emailRegex.MatchString(notification.Recipient) {
			return fmt.Errorf("invalid email address format: %s", notification.Recipient)
		}
	case model.ChannelSMS:
		if !phoneRegex.MatchString(notification.Recipient) {
			return fmt.Errorf("invalid phone number format: %s. Must be in E.164 format", notification.Recipient)
		}

		// Check message length (SMS has character limits)
		if len(notification.Message) > 160 {
			return fmt.Errorf("SMS message exceeds 160 character limit")
		}
	}

	return nil
}