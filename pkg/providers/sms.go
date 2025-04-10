package providers

import (
	"fmt"
	"log"
	"notification-system/pkg/config"
	"notification-system/pkg/model"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
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

func (t *TwilioSMSProvider) Send(notification model.Notification) error {
	params := &api.CreateMessageParams{}
	params.SetTo(notification.Recipient)
	params.SetFrom(t.fromNumber)
	params.SetBody(notification.Message)

	_, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	log.Printf("SMS sent successfully: %s", notification.Message)
	return nil
} 