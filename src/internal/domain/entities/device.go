package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserDevice struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	PushToken  string    `gorm:"unique;comment:Generic token for FCM/APNs" json:"push_token"`
	Provider   string    `gorm:"comment:fcm, apns" json:"provider"`            // fcm, apns
	DeviceType string    `gorm:"comment:ios, android, web" json:"device_type"` // ios, android, web
	LastSeen   time.Time `gorm:"default:now()" json:"last_seen"`
}
