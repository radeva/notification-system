package providers

import (
	"context"
	"fmt"
	"log"
	"notification-system/pkg/config"
	"notification-system/pkg/model"

	"github.com/slack-go/slack"
)

type SlackNotificationProvider struct {
	client *slack.Client
}

func NewSlackNotificationProvider(cfg config.SlackConfig) *SlackNotificationProvider {
	client := slack.New(cfg.BotToken)
	return &SlackNotificationProvider{
		client: client,
	}
}

func (s *SlackNotificationProvider) Send(ctx context.Context, notification model.Notification) error {
	// Create a channel to receive the result
	result := make(chan error, 1)

	// Start the send operation in a goroutine
	go func() {
		// Send the message to the specified channel
		_, _, err := s.client.PostMessageContext(
			ctx,
			notification.Recipient, // In Slack, recipient is the channel ID
			slack.MsgOptionText(notification.Message, false),
		)
		select {
		case result <- err:
			// Successfully sent result
		case <-ctx.Done():
			// Context was cancelled, discard the result
		}
	}()

	// Wait for either the operation to complete or the context to be done
	select {
	case err := <-result:
		if err != nil {
			return fmt.Errorf("failed to send Slack message: %w", err)
		}
		log.Printf("Slack message sent successfully: %s", notification.Message)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("slack message send operation cancelled: %w", ctx.Err())
	}
}