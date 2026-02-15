package entities

import (
	"github.com/google/uuid"
)

type UserSocialAccount struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Provider   string    `gorm:"comment:line, google" json:"provider"` // line, google
	ProviderID string    `gorm:"unique;comment:ID from social provider" json:"provider_id"`
}
