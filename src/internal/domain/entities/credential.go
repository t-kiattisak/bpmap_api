package entities

import (
	"github.com/google/uuid"
)

type SpecialCredential struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Username     string    `gorm:"unique" json:"username"`
	PasswordHash string    `json:"-"`
}
