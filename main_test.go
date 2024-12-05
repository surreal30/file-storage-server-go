package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPing(t *testing.T) {
	// Create a mock server that calls getPing
	mockServer := httptest.NewServer(http.HandlerFunc(getPing))
	defer mockServer.Close()

	// Send a GET request to the mock server
	resp, err := http.Get(mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
	assert.Equal(t, "pong working!", string(body), "Expected response body to be 'pong working!'")
}
