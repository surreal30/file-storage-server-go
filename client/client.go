package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
)

// Function to send request to the existing server
func pingServer() {
    serverURL := "http://localhost:2021/ping"

    // Send a GET request
    resp, err := http.Get(serverURL)
    if err != nil {
        log.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

    // Check if status code is OK
    if resp.StatusCode != http.StatusOK {
        log.Printf("Error: Received non-OK response status %d\n", resp.StatusCode)
        return
    }

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return
    }

    // Print the response body
    fmt.Printf("%s\n", string(body))
}

func main() {
    // Start listening for input commands from the user
    fmt.Println("CLI Program started. Type 'store' to send a request to the server.")

    // Create a scanner to read user input
    scanner := bufio.NewScanner(os.Stdin)

    for {
        // Prompt the user for a command
        fmt.Print("> ")
        scanner.Scan()

        // Get the user input and trim leading/trailing whitespace
        command := strings.TrimSpace(scanner.Text())

        // If 'store' is entered, send a request to the server
        if command == "store" {
            fmt.Println("Sending request to the server...")
            pingServer()
        } else if command == "exit" {
            // Exit the program if 'exit' is entered
            fmt.Println("Exiting program...")
            break
        } else {
            // Inform the user if the command is unknown
            fmt.Println("Unknown command:", command)
        }
    }
}
