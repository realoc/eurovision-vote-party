package middleware

import "context"

// WithUserID returns a context with the given user ID for testing purposes.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// WithUserEmail returns a context with the given email for testing purposes.
func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailContextKey, email)
}
