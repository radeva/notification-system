package services

import (
	"log"

	"notification-system/pkg/config"
	"notification-system/pkg/model"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSService handles sending SMS notifications
type SMSService struct {
	client *twilio.RestClient
	from   string
	enabled bool
}

// NewSMSService creates a new instance of SMSService
func NewSMSService(cfg config.TwilioConfig) *SMSService {
	// Get Twilio credentials from environment variables
	accountSid := cfg.AccountSID
	authToken := cfg.AuthToken
	fromNumber := cfg.FromNumber

	// Check if all required environment variables are set
	enabled := accountSid != "" && authToken != "" && fromNumber != ""
	
	if !enabled {
		log.Println("SMS service disabled: Missing Twilio credentials in environment variables")
		return &SMSService{
			enabled: false,
		}
	}

	// Create Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &SMSService{
		client: client,
		from:   fromNumber,
		enabled: true,
	}
}

// SendNotification sends an SMS
func (s *SMSService) SendNotification(n model.Notification) error {
	// If SMS service is disabled, just log and return
	if !s.enabled {
		log.Printf("SMS notification skipped (service disabled)")
		return nil
	}

	// Send SMS
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(n.Recipient)
	params.SetFrom(s.from)
	params.SetBody(n.Message)

	// Send the message
	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending SMS: %v", err)
		return err
	}

	log.Printf("SMS sent successfully: %s", n.Message)
	return nil
} 