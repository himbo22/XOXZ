package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	RoleID       uuid.UUID      `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID uuid.UUID      `gorm:"type:uuid;primaryKey" json:"permission_id"`
	GrantedBy    *uuid.UUID     `gorm:"type:uuid" json:"granted_by,omitempty"`
	Scope        *string        `json:"scope,omitempty"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Role          Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission    Permission `gorm:"foreignKey:PermissionID" json:"permission"`
	GrantedByUser User       `gorm:"foreignKey:GrantedBy" json:"granted_by_user"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
