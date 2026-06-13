package bootstrap

import (
	"sync"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
)

type Role struct {
	ID   uuid.UUID
	Name string
}

type Permission struct {
	ID   uuid.UUID
	Code string
}

type AccessControl struct {
	mu sync.RWMutex
	m  map[uuid.UUID]map[string]struct{}
}

func InitAccessControl(
	rp []entity.RolePermission,
	permissions []entity.Permission,
) (*AccessControl, func()) {
	permissionMap := make(map[uuid.UUID]string)

	for _, p := range permissions {
		permissionMap[p.ID] = p.Code
	}

	cache := make(map[uuid.UUID]map[string]struct{})

	for _, rolePermission := range rp {
		if _, ok := cache[rolePermission.RoleID]; !ok {
			cache[rolePermission.RoleID] = make(map[string]struct{})
		}

		code, ok := permissionMap[rolePermission.PermissionID]
		if !ok {
			continue
		}

		cache[rolePermission.RoleID][code] = struct{}{}
	}

	return &AccessControl{
		m: cache,
	}, nil
}

func (a *AccessControl) HasPermission(
	roleID uuid.UUID,
	code string,
) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	permissions, ok := a.m[roleID]
	if !ok {
		return false
	}

	_, ok = permissions[code]
	return ok
}
