package providers

import (
	"context"
	"errors"
	"notification-system/pkg/config"
	"notification-system/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// mockEmailClient implements the EmailClient interface for testing
type mockEmailClient struct {
	sendError error
	lastEmail *mail.SGMailV3
}

func (m *mockEmailClient) Send(email *mail.SGMailV3) error {
	return m.SendWithContext(context.Background(), email)
}

func (m *mockEmailClient) SendWithContext(ctx context.Context, email *mail.SGMailV3) error {
	m.lastEmail = email
	if m.sendError != nil {
		return m.sendError
	}
	return nil
}

var _ = Describe("EmailNotificationProvider", func() {
	var (
		provider     *EmailNotificationProvider
		mockClient   *mockEmailClient
		defaultCfg   *config.EmailConfig
		ctx          context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		defaultCfg = &config.EmailConfig{
			SendGridAPIKey: "test-key",
			FromAddress:    "test@example.com",
			FromName:       "Test Sender",
			DefaultSubject: "Test Subject",
		}
		mockClient = &mockEmailClient{}
		provider = &EmailNotificationProvider{
			config: defaultCfg,
			client: mockClient,
		}
	})

	Describe("Send", func() {
		Context("when sending a valid email", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "recipient@example.com",
					Message:   "Test content",
				}
			})

			It("should send the email successfully", func() {
				err := provider.Send(ctx, notification)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClient.lastEmail).NotTo(BeNil())
			})

			It("should use the default subject", func() {
				err := provider.Send(ctx, notification)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClient.lastEmail.Subject).To(Equal("Test Subject"))
			})

			It("should include the correct content", func() {
				err := provider.Send(ctx, notification)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClient.lastEmail.Content[0].Value).To(Equal("Test content"))
			})
		})

		Context("when sending an email with custom subject", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "recipient@example.com",
					Message:   "Test content",
					Metadata: map[string]string{
						"email_subject": "Custom Subject",
					},
				}
			})

			It("should use the custom subject from metadata", func() {
				err := provider.Send(ctx, notification)
				Expect(err).NotTo(HaveOccurred())
				Expect(mockClient.lastEmail.Subject).To(Equal("Custom Subject"))
			})
		})

		Context("when sending fails", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "recipient@example.com",
					Message:   "Test content",
				}
				mockClient.sendError = errors.New("sendgrid error")
			})

			It("should return the error", func() {
				err := provider.Send(ctx, notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sendgrid error"))
			})
		})

		Context("when context is cancelled", func() {
			var (
				notification model.Notification
				cancelCtx    context.Context
				cancel       context.CancelFunc
			)

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "recipient@example.com",
					Message:   "Test content",
				}
				cancelCtx, cancel = context.WithCancel(context.Background())
				cancel() // Cancel immediately
			})

			It("should return context cancellation error", func() {
				err := provider.Send(cancelCtx, notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cancelled"))
			})
		})
	})

	Describe("Validate", func() {
		Context("when notification is valid", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "recipient@example.com",
					Message:   "Test content",
				}
			})

			It("should not return an error", func() {
				err := provider.Validate(notification)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when notification has empty message", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "recipient@example.com",
					Message:   "",
				}
			})

			It("should return an error", func() {
				err := provider.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("message cannot be empty"))
			})
		})

		Context("when notification has empty recipient", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "",
					Message:   "Test content",
				}
			})

			It("should return an error", func() {
				err := provider.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recipient cannot be empty"))
			})
		})

		Context("when notification has invalid email format", func() {
			var notification model.Notification

			BeforeEach(func() {
				notification = model.Notification{
					ID:        "test-id",
					Channel:   model.ChannelEmail,
					Recipient: "invalid-email",
					Message:   "Test content",
				}
			})

			It("should return an error", func() {
				err := provider.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid email address format"))
			})
		})
	})
}) 