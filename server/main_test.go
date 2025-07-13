package main

import (
	"encoding/json"
	handlers "eurovision-app/http-handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupTestServer creates a test server with isolated routes
// and returns a function to create a new test request
func setupTestServer() func(method, path string, body string) *httptest.ResponseRecorder {
	// Create a new ServeMux for this test
	mux := http.NewServeMux()

	// Setup routes using the handlers package
	handlers.SetupRoutes(mux)

	// Return a function that creates a new test request
	return func(method, path string, body string) *httptest.ResponseRecorder {
		var req *http.Request
		var err error

		if body != "" {
			req, err = http.NewRequest(method, path, strings.NewReader(body))
			if err != nil {
				panic(err)
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(method, path, nil)
			if err != nil {
				panic(err)
			}
		}

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}
}

func TestHealthEndpoint(t *testing.T) {
	makeRequest := setupTestServer()
	rr := makeRequest("GET", "/health", "")

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Check the response body
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if status, exists := response["status"]; !exists || status != "healthy" {
		t.Errorf("handler returned unexpected body: got %v want status field with value 'healthy'",
			rr.Body.String())
	}
}

func TestPartyEndpoint(t *testing.T) {
	makeRequest := setupTestServer()
	payload := `{"party_name": "Test Party"}`
	rr := makeRequest("POST", "/party/create", payload)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Check the response body contains id and password
	var response handlers.PartyResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if response.ID == "" {
		t.Errorf("Response missing ID field")
	}

	if response.Password == "" {
		t.Errorf("Response missing Password field")
	}

	// Verify ID and Password have the expected length
	if len(response.ID) != 8 {
		t.Errorf("ID has unexpected length: got %d want 8", len(response.ID))
	}

	if len(response.Password) != 16 {
		t.Errorf("Password has unexpected length: got %d want 16", len(response.Password))
	}
}
