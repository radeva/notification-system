package providers

import (
	"context"
	"fmt"
	"log"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"regexp"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`) // E.164 format
)

type TwilioSMSProvider struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioSMSProvider(cfg config.TwilioConfig) *TwilioSMSProvider {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})

	return &TwilioSMSProvider{
		client:     client,
		fromNumber: cfg.FromNumber,
	}
}

func (t *TwilioSMSProvider) Send(ctx context.Context, notification model.Notification) error {
	// Create a channel to receive the result
	result := make(chan error, 1)

	// Start the send operation in a goroutine
	go func() {
		params := &api.CreateMessageParams{}
		params.SetTo(notification.Recipient)
		params.SetFrom(t.fromNumber)
		params.SetBody(notification.Message)

		_, err := t.client.Api.CreateMessage(params)
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
			return fmt.Errorf("failed to send SMS: %w", err)
		}
		log.Printf("SMS sent successfully: %s", notification.Message)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("SMS send operation timed out: %w", ctx.Err())
	}
}

func (s *TwilioSMSProvider) Validate(notification model.Notification) error {
	if notification.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	if notification.Recipient == "" {
		return fmt.Errorf("recipient cannot be empty")
	}

	if !phoneRegex.MatchString(notification.Recipient) {
		return fmt.Errorf("invalid phone number format: %s. Must be in E.164 format", notification.Recipient)
	}

	// Check message length (SMS has character limits)
	if len(notification.Message) > 160 {
		return fmt.Errorf("SMS message exceeds 160 character limit")
	}

	return nil
}