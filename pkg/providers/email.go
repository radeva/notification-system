package providers

import (
	"context"
	"fmt"
	"log"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"regexp"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	defaultEmailSubject = "Notification"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// EmailClient defines the interface for sending emails
type EmailClient interface {
	Send(email *mail.SGMailV3) error
	SendWithContext(ctx context.Context, email *mail.SGMailV3) error
}

// SendGridClient wraps the SendGrid client to implement our EmailClient interface
type SendGridClient struct {
	client *sendgrid.Client
}

func NewSendGridClient(apiKey string) *SendGridClient {
	return &SendGridClient{
		client: sendgrid.NewSendClient(apiKey),
	}
}

func (s *SendGridClient) Send(email *mail.SGMailV3) error {
	return s.SendWithContext(context.Background(), email)
}

func (s *SendGridClient) SendWithContext(ctx context.Context, email *mail.SGMailV3) error {
	response, err := s.client.SendWithContext(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("sendgrid API error: %d - %s", response.StatusCode, response.Body)
	}

	return nil
}

type EmailNotificationProvider struct {
	config *config.EmailConfig
	client EmailClient
}

func NewEmailNotificationProvider(cfg config.EmailConfig) *EmailNotificationProvider {
	return &EmailNotificationProvider{
		config: &cfg,
		client: NewSendGridClient(cfg.SendGridAPIKey),
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

		err := e.client.SendWithContext(ctx, message)
		if err != nil {
			result <- err
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

func (e *EmailNotificationProvider) Validate(notification model.Notification) error {
	if notification.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	if notification.Recipient == "" {
		return fmt.Errorf("recipient cannot be empty")
	}

	if !emailRegex.MatchString(notification.Recipient) {
		return fmt.Errorf("invalid email address format: %s", notification.Recipient)
	}

	return nil
} 