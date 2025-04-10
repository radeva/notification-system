package providers

import (
	"fmt"
	"notification-system/pkg/model"
	"sync"
)

// MockSMSProvider implements SMSProvider for testing
type MockSMSProvider struct {
	mu       sync.Mutex
	sent     []model.Notification
	FailNext bool
}

func NewMockSMSProvider() *MockSMSProvider {
	return &MockSMSProvider{
		sent: make([]model.Notification, 0),
	}
}

func (m *MockSMSProvider) Send(notification model.Notification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailNext {
		m.FailNext = false
		return fmt.Errorf("mock SMS provider failure")
	}

	m.sent = append(m.sent, notification)
	return nil
}

func (m *MockSMSProvider) GetSent() []model.Notification {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sent
}

func (m *MockSMSProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = make([]model.Notification, 0)
	m.FailNext = false
}

// MockEmailProvider implements EmailProvider for testing
type MockEmailProvider struct {
	mu       sync.Mutex
	sent     []model.Notification
	FailNext bool
}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{
		sent: make([]model.Notification, 0),
	}
}

func (m *MockEmailProvider) Send(notification model.Notification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailNext {
		m.FailNext = false
		return fmt.Errorf("mock email provider failure")
	}

	m.sent = append(m.sent, notification)
	return nil
}

func (m *MockEmailProvider) GetSent() []model.Notification {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sent
}

func (m *MockEmailProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = make([]model.Notification, 0)
	m.FailNext = false
}

// MockSlackProvider implements SlackProvider for testing
type MockSlackProvider struct {
	mu       sync.Mutex
	sent     []model.Notification
	FailNext bool
}

func NewMockSlackProvider() *MockSlackProvider {
	return &MockSlackProvider{
		sent: make([]model.Notification, 0),
	}
}

func (m *MockSlackProvider) Send(notification model.Notification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailNext {
		m.FailNext = false
		return fmt.Errorf("mock Slack provider failure")
	}

	m.sent = append(m.sent, notification)
	return nil
}

func (m *MockSlackProvider) GetSent() []model.Notification {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sent
}

func (m *MockSlackProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = make([]model.Notification, 0)
	m.FailNext = false
} 