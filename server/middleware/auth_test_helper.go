package middleware

import "context"

// WithUserID returns a context with the given user ID for testing purposes.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}
