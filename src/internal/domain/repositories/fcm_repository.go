package repositories

import (
	"context"
	"pbmap_api/src/internal/dto"
)

type FCMRepository interface {
	BroadcastNotification(ctx context.Context, title, body string) error
	SendAlarm(ctx context.Context, req *dto.AlarmDispatchRequest) error
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*dto.TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*dto.TopicManagementResponse, error)
}
