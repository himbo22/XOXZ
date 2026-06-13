package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"primaryKey" json:"id"`
	Username    *string        `json:"username"`
	Email       *string        `json:"email"`
	FirstName   *string        `json:"first_name"`
	LastName    *string        `json:"last_name"`
	Phone       *string        `json:"phone"`
	AvatarURL   *string        `json:"avatar_url"`
	Bio         *string        `json:"bio"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	Status      *string        `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}
