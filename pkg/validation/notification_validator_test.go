package validation

import (
	"notification-system/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationValidator", func() {
	var validator Validator

	BeforeEach(func() {
		validator = NewNotificationValidator()
	})

	Describe("Email notifications", func() {
		Context("with valid email", func() {
			It("should validate successfully", func() {
				notification := &model.Notification{
					Channel:   model.ChannelEmail,
					Recipient: "test@example.com",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid email", func() {
			It("should return an error for invalid format", func() {
				notification := &model.Notification{
					Channel:   model.ChannelEmail,
					Recipient: "invalid-email",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid email address format"))
			})
		})
	})

	Describe("SMS notifications", func() {
		Context("with valid phone number", func() {
			It("should validate successfully", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSMS,
					Recipient: "+1234567890",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid phone number", func() {
			It("should return an error for invalid format", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSMS,
					Recipient: "1234567890", // Missing + prefix
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid phone number format"))
			})
		})

		Context("with message length", func() {
			It("should return an error for message exceeding 160 characters", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSMS,
					Recipient: "+1234567890",
					Message:   "This is a very long message that exceeds the 160 character limit for SMS messages. This is a very long message that exceeds the 160 character limit for SMS messages. This is a very long message that exceeds the 160 character limit for SMS messages.",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SMS message exceeds 160 character limit"))
			})
		})
	})

	Describe("Slack notifications", func() {
		Context("with valid channel ID", func() {
			It("should validate successfully for channel ID", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSlack,
					Recipient: "C123456789",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should validate successfully for group ID", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSlack,
					Recipient: "G123456789",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid channel ID", func() {
			It("should return an error for wrong prefix", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSlack,
					Recipient: "X123456789",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid slack channel ID format"))
			})

			It("should return an error for ID too short", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSlack,
					Recipient: "C1234567",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid slack channel ID format"))
			})

			It("should return an error for ID too long", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSlack,
					Recipient: "C12345678901",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid slack channel ID format"))
			})

			It("should return an error for invalid characters", func() {
				notification := &model.Notification{
					Channel:   model.ChannelSlack,
					Recipient: "C1234567-9",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid slack channel ID format"))
			})
		})
	})

	Describe("Common validation", func() {
		Context("with empty message", func() {
			It("should return an error", func() {
				notification := &model.Notification{
					Channel:   model.ChannelEmail,
					Recipient: "test@example.com",
					Message:   "",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("message cannot be empty"))
			})
		})

		Context("with empty recipient", func() {
			It("should return an error", func() {
				notification := &model.Notification{
					Channel:   model.ChannelEmail,
					Recipient: "",
					Message:   "Test message",
				}
				err := validator.Validate(notification)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recipient cannot be empty"))
			})
		})
	})
}) 