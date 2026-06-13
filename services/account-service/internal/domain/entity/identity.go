package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Identity struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID      `json:"user_id"`
	Provider       string         `json:"provider"`
	ProviderUserID string         `json:"provider_user_id"`
	ProviderData   datatypes.JSON `json:"provider_data"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

func (Identity) TableName() string {
	return "identities"
}
