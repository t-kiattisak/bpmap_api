package usecase

import (
	"context"

	"pbmap_api/src/internal/domain"
)

// AlarmUsecase orchestrates alarm dispatch.
type AlarmUsecase interface {
	DispatchAlarm(ctx context.Context, req *domain.AlarmDispatchRequest) error
}

type alarmUsecase struct {
	fcm domain.FCMService
}

// NewAlarmUsecase creates the alarm usecase.
func NewAlarmUsecase(fcm domain.FCMService) AlarmUsecase {
	return &alarmUsecase{fcm: fcm}
}

func (u *alarmUsecase) DispatchAlarm(ctx context.Context, req *domain.AlarmDispatchRequest) error {
	return u.fcm.SendAlarm(ctx, req)
}
