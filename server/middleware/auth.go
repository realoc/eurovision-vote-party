package middleware

import (
	"context"
	"net/http"
	"strings"

	firebaseauth "firebase.google.com/go/v4/auth"
)

// tokenVerifier defines the subset of the Firebase Auth client used by the middleware.
type tokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*firebaseauth.Token, error)
}

// contextKey avoids collisions with other context values.
type contextKey string

const (
	// userIDContextKey stores the Firebase user ID extracted from a verified token.
	userIDContextKey contextKey = "middleware/firebaseUserID"
	// userEmailContextKey stores the email extracted from a verified token.
	userEmailContextKey contextKey = "middleware/firebaseUserEmail"
)

var verifier tokenVerifier

// SetTokenVerifier configures the package-level verifier used by the auth middleware.
// It must be called during application startup before AuthMiddleware is used.
func SetTokenVerifier(v tokenVerifier) {
	verifier = v
}

// AuthMiddleware verifies Firebase ID tokens from the Authorization header.
// For valid tokens, the Firebase user ID is attached to the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if verifier == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		verifiedToken, err := verifier.VerifyIDToken(r.Context(), token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, verifiedToken.UID)
		if email, ok := verifiedToken.Claims["email"].(string); ok {
			ctx = context.WithValue(ctx, userEmailContextKey, email)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext extracts the Firebase user ID from the request context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(userIDContextKey).(string)
	return uid, ok
}

// UserEmailFromContext extracts the email from the request context.
func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailContextKey).(string)
	return email, ok
}

// OptionalAuthMiddleware extracts the user ID from the Authorization header if present.
// Unlike AuthMiddleware, it does not block requests without valid auth - it simply
// passes them through without a user ID in the context.
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if verifier == nil {
			next.ServeHTTP(w, r)
			return
		}

		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			next.ServeHTTP(w, r)
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		verifiedToken, err := verifier.VerifyIDToken(r.Context(), token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, verifiedToken.UID)
		if email, ok := verifiedToken.Claims["email"].(string); ok {
			ctx = context.WithValue(ctx, userEmailContextKey, email)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
