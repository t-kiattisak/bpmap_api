package usecase

import (
	"context"
	"fmt"
)

type DataSyncService interface {
	SyncAndNotify(ctx context.Context) error
}

type dataSyncService struct {
	fcmService FCMService
}

func NewDataSyncService(fcmService FCMService) DataSyncService {
	return &dataSyncService{
		fcmService: fcmService,
	}
}

func (s *dataSyncService) SyncAndNotify(ctx context.Context) error {
	fmt.Println("[DataSync] Starting sync job...")

	data, err := s.fetchExternalData()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}

	result := s.calculateLogic(data)

	if result.ShouldNotify {
		err := s.fcmService.BroadcastNotification(ctx, "New Data Update", result.Message)
		if err != nil {
			return fmt.Errorf("failed to broadcast notification: %v", err)
		}
	}

	fmt.Println("[DataSync] Job completed successfully.")
	return nil
}

type ExternalDataStub struct {
	Value int
}

type CalculationResult struct {
	ShouldNotify bool
	Message      string
}

func (s *dataSyncService) fetchExternalData() (*ExternalDataStub, error) {
	fmt.Println("[DataSync] Fetching external api...")
	return &ExternalDataStub{Value: 100}, nil
}

func (s *dataSyncService) calculateLogic(data *ExternalDataStub) *CalculationResult {
	fmt.Println("[DataSync] Calculating logic...")

	if data.Value > 50 {
		return &CalculationResult{
			ShouldNotify: true,
			Message:      fmt.Sprintf("Alert! Value is high: %d", data.Value),
		}
	}

	return &CalculationResult{ShouldNotify: false}
}
