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

		ctxWithUID := context.WithValue(r.Context(), userIDContextKey, verifiedToken.UID)
		next.ServeHTTP(w, r.WithContext(ctxWithUID))
	})
}

// UserIDFromContext extracts the Firebase user ID from the request context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(userIDContextKey).(string)
	return uid, ok
}
