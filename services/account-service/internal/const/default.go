package _const

type contextKey string

const TxKey contextKey = "tx"

type headerKey string

const HeaderDeviceID headerKey = "Device-ID"
const HeaderDeviceType headerKey = "Device-Type"

type DeviceType string

const (
	WebType    DeviceType = "1"
	MobileType DeviceType = "2"
	NilType    DeviceType = "0"
)

var DeviceTypeAllowed = map[DeviceType]struct{}{
	WebType:    {}, // web
	MobileType: {}, // mobile
}

// role
const (
	ROLE_ADMIN = "ADMIN"
	ROLE_IDOL  = "IDOL"
	ROLE_USER  = "USER"
)

// permission
const ()

const (
	USER_STATUS_ACTIVE = "ACTIVE"
)

// ContextKey type for context keys
type ContextKey string

const (
	// IPAddressKey is the context key for IP address
	IPAddressKey ContextKey = "ipa"
	// TimezoneKey is the context key for timezone
	TimezoneKey ContextKey = "tz"
	// UserAgentKey is the context key for user agent
	UserAgentKey ContextKey = "ua"
	// Device type
	DeviceTypeKey ContextKey = "dtk"
	// Token
	AccessTokenKey        ContextKey = "at"
	RefreshTokenCookieKey ContextKey = "rt"
	// User
	UserIDKey   ContextKey = "uid"
	RoleIDKey   ContextKey = "rid"
	DeviceIDKey ContextKey = "did"
	IssuerAtKey ContextKey = "iat"
)
