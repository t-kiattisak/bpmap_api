package usecase

import (
	"context"

	"pbmap_api/src/internal/domain"
)

// NotificationUsecase orchestrates notification (broadcast, subscribe, unsubscribe).
type NotificationUsecase interface {
	Broadcast(ctx context.Context, title, body string) error
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*domain.TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*domain.TopicManagementResponse, error)
}

type notificationUsecase struct {
	fcm domain.FCMService
}

// NewNotificationUsecase creates the notification usecase.
func NewNotificationUsecase(fcm domain.FCMService) NotificationUsecase {
	return &notificationUsecase{fcm: fcm}
}

func (u *notificationUsecase) Broadcast(ctx context.Context, title, body string) error {
	return u.fcm.BroadcastNotification(ctx, title, body)
}

func (u *notificationUsecase) SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*domain.TopicManagementResponse, error) {
	return u.fcm.SubscribeToTopic(ctx, tokens, topic)
}

func (u *notificationUsecase) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*domain.TopicManagementResponse, error) {
	return u.fcm.UnsubscribeFromTopic(ctx, tokens, topic)
}
