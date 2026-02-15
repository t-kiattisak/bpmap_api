package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `gorm:"unique" json:"refresh_token"`
	DeviceID     uuid.UUID `json:"device_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
}
