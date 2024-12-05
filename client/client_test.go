package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func testGetFiles(t *testing.T) {
	// mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/list" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "File ID: 12, Name: abc.txt\nFile ID: 13, Name: file1.txt\nFile ID: 14, Name: temp.txt")
	}))
	defer mockServer.Close()

	serverURL := mockServer.URL + "/list"

	result := getFiles(serverURL)
	if result != "File ID: 12, Name: abc.txt\nFile ID: 13, Name: file1.txt\nFile ID: 14, Name: temp.txt" {
		t.Errorf("Incorrect output")
	}
}

func TestDeleteFile(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the method is DELETE and the URL is correct
		if r.Method != http.MethodDelete || r.URL.Path != "/delete" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// Check if the file part exists in the form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}
		file, _, err := r.FormFile("files")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Simulate successful file deletion
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "File successfully deleted!")
	}))
	defer mockServer.Close()

	// Replace the baseURL in the function with the mock server's URL
	baseURL := mockServer.URL

	// Create a temporary test file for testing
	testFileName := "testfile.txt"
	err := ioutil.WriteFile(testFileName, []byte("This is a file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFileName) // Cleanup after test

	// Call the deleteFile function
	result, err := deleteFile(baseURL, testFileName)

	// Check that no error occurred and the result is as expected
	assert.NoError(t, err, "Expected no error when deleting file")
	assert.Equal(t, "File successfully deleted!", result, "Expected file deletion success message")
}

func TestPostFile(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the method is POST and the URL is correct
		if r.Method != http.MethodPost || r.URL.Path != "/add" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			http.Error(w, "No files found in the form", http.StatusInternalServerError)
			return
		}

		for _, fileHeader := range files {
			fmt.Println("Received file:", fileHeader.Filename)
		}

		// Simulate successful file upload
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Files created successfully!")
	}))
	defer mockServer.Close()

	testFileName := "testfile.txt"
	err := ioutil.WriteFile(testFileName, []byte("This is a file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFileName) // Cleanup after test

	result, err := postFile(mockServer.URL, []string{testFileName})

	assert.NoError(t, err, "Expected no error when creating file")
	assert.Equal(t, "Files created successfully!", result, "Expected file creation success message success message")
}

func TestPutFile(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/update" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}
		file, _, err := r.FormFile("files")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "File updated successfully!")
	}))
	defer mockServer.Close()

	testFileName := "testfile.txt"
	err := ioutil.WriteFile(testFileName, []byte("This is a file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFileName)


	result, err := putFile(mockServer.URL, testFileName)

	assert.NoError(t, err, "Expected no error when updating file")
	assert.Equal(t, "File updated successfully!\n", result, "Expected file deletion success message")
}