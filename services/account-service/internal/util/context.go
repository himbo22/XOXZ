package util

import "context"

type contextKey string

const UserIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(UserIDKey)
	if id, ok := val.(string); ok {
		return id, true
	}
	return "", false
}
