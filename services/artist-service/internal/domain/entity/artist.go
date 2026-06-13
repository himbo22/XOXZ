package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Artist struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	StageName   string         `gorm:"size:120;not null;uniqueIndex" json:"stage_name"`
	DisplayName *string        `gorm:"size:120" json:"display_name"`
	Bio         *string        `gorm:"type:text" json:"bio"`
	AvatarURL   *string        `gorm:"type:text" json:"avatar_url"`
	BannerURL   *string        `gorm:"type:text" json:"banner_url"`
	Verified    bool           `gorm:"not null;default:false" json:"verified"`
	Status      string         `gorm:"size:32;not null;default:'active'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Artist) TableName() string {
	return "artists"
}
