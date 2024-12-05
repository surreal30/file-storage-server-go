package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testGetWC(t *testing.T) {
	// mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wc" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "All files contain 33 words")
	}))
	defer mockServer.Close()

	serverURL := mockServer.URL + "/wc"

	result, err := getWC(serverURL)
	if result != "All files contain 33 words" && err == nil {
		t.Errorf("Incorrect output")
	}
}


func testPingServer(t *testing.T) {
	// mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong working!")
	}))
	defer mockServer.Close()

	serverURL := mockServer.URL + "/ping"

	result := pingServer(serverURL)
	if result != "pong working!" {
		t.Errorf("Incorrect output")
	}
}