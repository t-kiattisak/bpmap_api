package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email       *string    `gorm:"unique" json:"email"`
	DisplayName string     `json:"display_name"`
	Role        string     `gorm:"type:varchar(20);comment:citizen, officer, admin" json:"role"` // citizen, officer, admin
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	SocialAccounts    []UserSocialAccount `gorm:"foreignKey:UserID" json:"social_accounts,omitempty"`
	SpecialCredential *SpecialCredential  `gorm:"foreignKey:UserID" json:"special_credential,omitempty"`
	Devices           []UserDevice        `gorm:"foreignKey:UserID" json:"devices,omitempty"`
	Sessions          []UserSession       `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}
