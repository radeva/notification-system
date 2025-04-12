package providers

import (
	"context"
	"fmt"
	"notification-system/pkg/model"
)

// NotificationStrategyContext manages the notification sending strategies
type NotificationStrategyContext struct {
	strategies map[model.NotificationChannel]NotificationProvider
}

func NewNotificationStrategyContext() *NotificationStrategyContext {
	return &NotificationStrategyContext{
		strategies: make(map[model.NotificationChannel]NotificationProvider),
	}
}

// RegisterStrategy registers a provider for a specific channel
func (c *NotificationStrategyContext) RegisterStrategy(channel model.NotificationChannel, provider NotificationProvider) {
	c.strategies[channel] = provider
}

// Send uses the appropriate strategy based on the notification channel
func (c *NotificationStrategyContext) Send(ctx context.Context, notification model.Notification) error {
	provider, exists := c.strategies[notification.Channel]
	if !exists {
		return fmt.Errorf("no provider registered for channel: %s", notification.Channel)
	}
	return provider.Send(ctx, notification)
}

func (n *NotificationStrategyContext) GetStrategy(channel model.NotificationChannel) (NotificationProvider, error) {
	provider, exists := n.strategies[channel]
	if !exists {
		return nil, fmt.Errorf("no provider registered for channel %s", channel)
	}
	return provider, nil
} 