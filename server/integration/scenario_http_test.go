//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	firebaseauth "firebase.google.com/go/v4/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/handlers"
	"github.com/sipgate/eurovision-vote-party/server/middleware"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// stubTokenVerifier implements the tokenVerifier interface via structural typing.
// It maps token strings to (uid, email) pairs for testing.
type stubTokenVerifier struct {
	tokens map[string]tokenInfo
}

type tokenInfo struct {
	uid   string
	email string
}

func (s *stubTokenVerifier) VerifyIDToken(_ context.Context, idToken string) (*firebaseauth.Token, error) {
	info, ok := s.tokens[idToken]
	if !ok {
		return nil, http.ErrAbortHandler // any non-nil error signals invalid token
	}
	return &firebaseauth.Token{
		UID: info.uid,
		Claims: map[string]interface{}{
			"email": info.email,
		},
	}, nil
}

// buildMux constructs the same http.ServeMux as main.go, using real services
// backed by the Firestore emulator.
func buildMux(t *testing.T) *http.ServeMux {
	t.Helper()

	partyDAO := persistence.NewFirestorePartyDAO(firestoreClient)
	guestDAO := persistence.NewFirestoreGuestDAO(firestoreClient)
	voteDAO := persistence.NewFirestoreVoteDAO(firestoreClient)
	userDAO := persistence.NewFirestoreUserDAO(firestoreClient)

	partyService := services.NewPartyService(partyDAO)
	guestService := services.NewGuestService(guestDAO, partyDAO)
	voteService := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
	userService := services.NewUserService(userDAO)

	partyHandler := handlers.NewPartyHandler(partyService)
	guestHandler := handlers.NewGuestHandler(guestService)
	voteHandler := handlers.NewVoteHandler(voteService)
	actsHandler := handlers.NewActsHandler(actsService)
	userHandler := handlers.NewUserHandler(userService)

	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/parties/")
		segments := strings.SplitN(path, "/", 3)
		if len(segments) >= 2 {
			switch segments[1] {
			case "votes", "end-voting", "results":
				voteHandler.ServeHTTP(w, r)
				return
			}
			guestHandler.ServeHTTP(w, r)
			return
		}
		partyHandler.ServeHTTP(w, r)
	})

	mux := http.NewServeMux()
	mux.Handle("/api/health", handlers.NewHealthHandler())
	mux.Handle("/api/acts", actsHandler)
	mux.Handle("/api/parties", middleware.AuthMiddleware(partyHandler))
	mux.Handle("/api/parties/", middleware.OptionalAuthMiddleware(apiHandler))
	mux.Handle("/api/users/profile", middleware.AuthMiddleware(userHandler))

	return mux
}

func TestHTTPFullStack(t *testing.T) {
	// Set up stub verifier
	stub := &stubTokenVerifier{
		tokens: map[string]tokenInfo{
			"admin-token": {uid: "http-admin-1", email: "admin@test.com"},
			"guest-token": {uid: "http-guest-1", email: "guest@test.com"},
		},
	}
	middleware.SetTokenVerifier(stub)

	t.Cleanup(func() {
		ctx := context.Background()
		for _, col := range []string{"parties", "guests", "votes", "users"} {
			cleanupCollection(t, ctx, col)
		}
	})

	mux := buildMux(t)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "ok", body["status"])
	})

	t.Run("list acts", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/acts")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Acts []models.Act `json:"acts"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.NotEmpty(t, body.Acts)
	})

	t.Run("list acts with event filter", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/acts?event=grandfinal")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Acts []models.Act `json:"acts"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		for _, act := range body.Acts {
			assert.Equal(t, models.EventGrandFinal, act.EventType)
		}
	})

	t.Run("create party requires auth", func(t *testing.T) {
		reqBody := `{"name":"No Auth Party","eventType":"grandfinal"}`
		resp, err := http.Post(server.URL+"/api/parties", "application/json", strings.NewReader(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	var partyCode string
	var partyID string

	t.Run("create party with auth", func(t *testing.T) {
		reqBody := `{"name":"HTTP Test Party","eventType":"grandfinal"}`
		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/parties", strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer admin-token")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var party models.Party
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&party))
		assert.Equal(t, "HTTP Test Party", party.Name)
		assert.NotEmpty(t, party.Code)
		assert.NotEmpty(t, party.ID)
		partyCode = party.Code
		partyID = party.ID
	})

	t.Run("guest joins party without auth", func(t *testing.T) {
		require.NotEmpty(t, partyCode, "party must be created first")

		reqBody := `{"username":"HTTPGuest"}`
		resp, err := http.Post(
			server.URL+"/api/parties/"+partyCode+"/join",
			"application/json",
			strings.NewReader(reqBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var guest models.Guest
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&guest))
		assert.Equal(t, "HTTPGuest", guest.Username)
		assert.Equal(t, models.GuestStatusPending, guest.Status)
	})

	t.Run("get party by code without auth", func(t *testing.T) {
		require.NotEmpty(t, partyCode, "party must be created first")

		resp, err := http.Get(server.URL + "/api/parties/" + partyCode)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, partyCode, body["code"])
		// Public endpoint should not expose adminId
		_, hasAdminID := body["adminId"]
		assert.False(t, hasAdminID, "public party response should not expose adminId")
	})

	t.Run("admin list join requests", func(t *testing.T) {
		require.NotEmpty(t, partyID, "party must be created first")

		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/parties/"+partyID+"/join-requests", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer admin-token")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var guests []models.Guest
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&guests))
		require.Len(t, guests, 1)
		assert.Equal(t, "HTTPGuest", guests[0].Username)
	})

	t.Run("user profile upsert and get", func(t *testing.T) {
		// PUT profile
		reqBody := `{"username":"admin_user"}`
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/users/profile", bytes.NewBufferString(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer admin-token")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var user models.User
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&user))
		assert.Equal(t, "admin_user", user.Username)
		assert.Equal(t, "admin@test.com", user.Email)

		// GET profile
		req2, err := http.NewRequest(http.MethodGet, server.URL+"/api/users/profile", nil)
		require.NoError(t, err)
		req2.Header.Set("Authorization", "Bearer admin-token")

		resp2, err := http.DefaultClient.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		var user2 models.User
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&user2))
		assert.Equal(t, "admin_user", user2.Username)
	})
}
