package util

import (
	_const "github.com/himbo22/xoxz/account-service/internal/const"
)

func IsValidDeviceType(deviceType _const.DeviceType) bool {
	_, ok := _const.DeviceTypeAllowed[deviceType]
	return ok
}

func ParseDeviceType(deviceType string) _const.DeviceType {
	switch deviceType {
	case "1":
		return _const.WebType
	case "2":
		return _const.MobileType
	default:
		return _const.NilType
	}
}
