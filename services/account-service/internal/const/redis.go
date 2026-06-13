package _const

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	refreshTokenKey         = "rt:%v"   // rt:{rt_value}	|	{userID , deviceID}
	userDevicesKey          = "u:d:%v"  // Set: u:d:{uid} | Members: Hash[{device_id}: {refresh_token}]
	blackAccessTokenListKey = "u:ra:%v" // u:ra:{userID} | timeline: 8:00
	googleTokenKey          = "g:%v"    // g:{google-token} | 1 (keep payload as light as possible)
)

func RefreshTokenKey(refreshToken string) string {
	return fmt.Sprintf(refreshTokenKey, refreshToken)
}

func UserDevicesKey(userID uuid.UUID) string {
	return fmt.Sprintf(userDevicesKey, userID)
}

func BlackAccessTokenKey(userID uuid.UUID) string {
	return fmt.Sprintf(blackAccessTokenListKey, userID)
}

func GoogleTokenKey(token string) string {
	return fmt.Sprintf(googleTokenKey, token)
}
