package util

import "time"

func UnixToDuration(ttl int64) time.Duration {
	return time.Until(time.Unix(ttl, 0))
}

func StrPtr(s string) *string {
	return &s
}
