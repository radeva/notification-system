package providers

import (
	"context"
	"fmt"
	"log"
	"notification-system/pkg/config"
	"notification-system/pkg/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	defaultEmailSubject = "Notification"
)

type EmailNotificationProvider struct {
	config *config.EmailConfig
	client *sendgrid.Client
}

func NewEmailNotificationProvider(cfg config.EmailConfig) *EmailNotificationProvider {
	client := sendgrid.NewSendClient(cfg.SendGridAPIKey)
	return &EmailNotificationProvider{
		config: &cfg,
		client: client,
	}
}

func (e *EmailNotificationProvider) Send(ctx context.Context, notification model.Notification) error {
	// Create a channel to receive the result
	result := make(chan error, 1)

	// Start the send operation in a goroutine
	go func() {
		from := mail.NewEmail(e.config.FromName, e.config.FromAddress)
		to := mail.NewEmail("", notification.Recipient)
		
		// Get subject from metadata for email notifications, fallback to configured default
		subject := e.config.DefaultSubject
		if subject == "" {
			subject = defaultEmailSubject // Fallback to hardcoded default if not configured
		}
		if notification.Metadata != nil {
			if emailSubject, ok := notification.Metadata["email_subject"]; ok {
				subject = emailSubject
			}
		}

		plainTextContent := notification.Message
		htmlContent := fmt.Sprintf("<p>%s</p>", notification.Message)

		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

		response, err := e.client.SendWithContext(ctx, message)
		if err != nil {
			result <- fmt.Errorf("failed to send email: %w", err)
			return
		}

		if response.StatusCode >= 300 {
			result <- fmt.Errorf("sendgrid API error: %d - %s", response.StatusCode, response.Body)
			return
		}

		select {
		case result <- nil:
			// Successfully sent result
		case <-ctx.Done():
			// Context was cancelled, discard the result
		}
	}()

	// Wait for either the operation to complete or the context to be done
	select {
	case err := <-result:
		if err != nil {
			return err
		}
		log.Printf("Email sent successfully to %s: %s", notification.Recipient, notification.Message)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("email send operation cancelled: %w", ctx.Err())
	}
} 