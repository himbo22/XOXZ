package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole struct {
	RoleID     uuid.UUID      `gorm:"type:uuid;primaryKey" json:"role_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;primaryKey" json:"user_id"`
	AssignedBy *uuid.UUID     `gorm:"type:uuid" json:"assigned_by,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Role           Role `gorm:"foreignKey:RoleID" json:"role"`
	User           User `gorm:"foreignKey:UserID" json:"user"`
	AssignedByUser User `gorm:"foreignKey:AssignedBy" json:"assigned_by_user"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
