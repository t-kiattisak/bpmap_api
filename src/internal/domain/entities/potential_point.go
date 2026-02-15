package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PotentialPoint struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name       string         `gorm:"type:varchar(255);not null"`
	Type       string         `gorm:"type:varchar(50);not null"`
	Latitude   float64        `gorm:"type:decimal(10,8);not null"`
	Longitude  float64        `gorm:"type:decimal(11,8);not null"`
	CreatedBy  uuid.UUID      `gorm:"type:uuid;not null"`
	Properties datatypes.JSON `gorm:"type:jsonb"`

	Creator   *User      `gorm:"foreignKey:CreatedBy"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}
