package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	firebaseauth "firebase.google.com/go/v4/auth"
)

type stubTokenVerifier struct {
	token      *firebaseauth.Token
	err        error
	callCount  int
	receivedID string
}

func (s *stubTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
	s.callCount++
	s.receivedID = idToken
	return s.token, s.err
}

func TestAuthMiddlewareRejectsMissingAuthorizationHeader(t *testing.T) {
	t.Cleanup(func() {
		SetTokenVerifier(nil)
	})

	SetTokenVerifier(&stubTokenVerifier{})

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	AuthMiddleware(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidAuthorizationHeaderFormat(t *testing.T) {
	t.Cleanup(func() {
		SetTokenVerifier(nil)
	})

	SetTokenVerifier(&stubTokenVerifier{})

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic token")
	rec := httptest.NewRecorder()

	AuthMiddleware(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAuthMiddlewareRejectsWhenVerifierReturnsError(t *testing.T) {
	t.Cleanup(func() {
		SetTokenVerifier(nil)
	})

	stub := &stubTokenVerifier{
		err: errors.New("invalid token"),
	}

	SetTokenVerifier(stub)

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	rec := httptest.NewRecorder()

	AuthMiddleware(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if stub.callCount != 1 {
		t.Fatalf("expected verifier to be called once, got %d", stub.callCount)
	}

	if stub.receivedID != "badtoken" {
		t.Fatalf("expected token \"badtoken\", got %q", stub.receivedID)
	}
}

func TestAuthMiddlewareAllowsRequestWhenTokenValid(t *testing.T) {
	t.Cleanup(func() {
		SetTokenVerifier(nil)
	})

	stub := &stubTokenVerifier{
		token: &firebaseauth.Token{
			UID: "user-123",
		},
	}

	SetTokenVerifier(stub)

	var (
		handlerCalled bool
		capturedUID   string
	)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		uid, ok := UserIDFromContext(r.Context())
		if !ok {
			t.Fatalf("expected user ID in context")
		}
		capturedUID = uid
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	rec := httptest.NewRecorder()

	AuthMiddleware(next).ServeHTTP(rec, req)

	if !handlerCalled {
		t.Fatalf("expected next handler to be called")
	}

	if capturedUID != "user-123" {
		t.Fatalf("expected UID \"user-123\", got %q", capturedUID)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if stub.callCount != 1 {
		t.Fatalf("expected verifier to be called once, got %d", stub.callCount)
	}

	if stub.receivedID != "validtoken" {
		t.Fatalf("expected token \"validtoken\", got %q", stub.receivedID)
	}
}

func TestAuthMiddlewareExtractsEmailFromToken(t *testing.T) {
	t.Cleanup(func() {
		SetTokenVerifier(nil)
	})

	stub := &stubTokenVerifier{
		token: &firebaseauth.Token{
			UID: "user-123",
			Claims: map[string]interface{}{
				"email": "user@example.com",
			},
		},
	}

	SetTokenVerifier(stub)

	var capturedEmail string
	var emailFound bool

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedEmail, emailFound = UserEmailFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	rec := httptest.NewRecorder()

	AuthMiddleware(next).ServeHTTP(rec, req)

	if !emailFound {
		t.Fatalf("expected email in context")
	}

	if capturedEmail != "user@example.com" {
		t.Fatalf("expected email \"user@example.com\", got %q", capturedEmail)
	}
}

func TestAuthMiddlewareHandsMissingEmailGracefully(t *testing.T) {
	t.Cleanup(func() {
		SetTokenVerifier(nil)
	})

	stub := &stubTokenVerifier{
		token: &firebaseauth.Token{
			UID: "user-123",
			// No email claim
		},
	}

	SetTokenVerifier(stub)

	var emailFound bool

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, emailFound = UserEmailFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	rec := httptest.NewRecorder()

	AuthMiddleware(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if emailFound {
		t.Fatalf("expected no email in context when token has no email claim")
	}
}
