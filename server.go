package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Simple function to ping and test if server is up or not
func getPing(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("pong!\n")
	io.WriteString(w, "pong working!")
}

func main() {
	http.HandleFunc("/ping", getPing)

	err := http.ListenAndServe(":2021", nil)

	if err != nil {
		fmt.Printf("Error starting in server: %s\n", err)
		os.Exit(1)
	}

}