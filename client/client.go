package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "log"
    "mime/multipart"
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

func getFiles() {
    serverURL := "http://localhost:2021/list"

    resp, err := http.Get(serverURL)
    if err != nil {
        log.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error: Received non-OK response status %d\n", resp.StatusCode)
        return
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return
    }

    fmt.Printf("%s\n", string(body))

}

func deleteFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("Error opening %s: %v", filename, err)
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Add the file to the request body
    part, err := writer.CreateFormFile("files", filename)
    if err != nil {
        return fmt.Errorf("Error creating form file: %v", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return fmt.Errorf("Error copying file content: %v", err)
    }

    err = writer.Close()
    if err != nil {
        return fmt.Errorf("Error closing writer: %v", err)
    }

    url := "http://localhost:2021/delete"
    req, err := http.NewRequest("DELETE", url, &requestBody)
    if err != nil {
        return fmt.Errorf("Error creating request: %v", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Error: received non-OK response: %v", resp.Status)
    }

    fmt.Println("File successfully deleted!")
    return nil
}

func putFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("Error opening %s: %v", filename, err)
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Add the file to the request body
    part, err := writer.CreateFormFile("files", filename)
    if err != nil {
        return fmt.Errorf("Error creating form file: %v", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return fmt.Errorf("Error copying file content: %v", err)
    }

    err = writer.Close()
    if err != nil {
        return fmt.Errorf("Error closing writer: %v", err)
    }

    url := "http://localhost:2021/update"
    req, err := http.NewRequest("PUT", url, &requestBody)
    if err != nil {
        return fmt.Errorf("Error creating request: %v", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Error: received non-OK response: %v", resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error updating the file: ", err)
        return nil
    }

    fmt.Printf("%s\n", string(body))
    return nil
}

func postFile(filenames []string) error {
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    for _, filename := range filenames {
        fmt.Println("Attempting to open file:", filename)

        file, err := os.Open(filename)
        if err != nil {
            return fmt.Errorf("Error opening file '%s': %v", filename, err)
        }
        defer file.Close()

        part, err := writer.CreateFormFile("files", filename)
        if err != nil {
            return fmt.Errorf("Error creating form file for '%s': %v", filename, err)
        }

        _, err = io.Copy(part, file)
        if err != nil {
            return fmt.Errorf("Error copying file content for '%s': %v", filename, err)
        }
    }

    err := writer.Close()
    if err != nil {
        return fmt.Errorf("Error closing writer: %v", err)
    }

    url := "http://localhost:2021/add"
    req, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        return fmt.Errorf("Error creating request: %v", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Error: received non-OK response: %v", resp.Status)
    }

    responseBody, _ := io.ReadAll(resp.Body)
    fmt.Println("Server Response:", string(responseBody))

    fmt.Println("Files created successfully!")
    return nil
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
        if strings.HasPrefix(command, "store rm ") {
            filename := strings.TrimPrefix(command, "store rm ")
            fmt.Printf("Sending delete request for file: %s\n", filename)
            err := deleteFile(filename)
            if err != nil {
                log.Printf("Error: %v\n", err)
            }
        } else if strings.HasPrefix(command, "store update ") {
            filename := strings.TrimPrefix(command, "store update ")
            fmt.Printf("Sending update request for file: %s\n", filename)
            err := putFile(filename)
            if err != nil {
                log.Printf("Error: %v\n", err)
            }
        } else if strings.HasPrefix(command, "store add") {
            parts := strings.Fields(command)
            filenames := parts[2:]
            fmt.Println("Sending create request\n")
            err := postFile(filenames)
            if err != nil {
                log.Printf("Error: %v\n", err)
            }
        } else if command == "store" {
            fmt.Println("Sending request to the server...")
            pingServer()
        } else if command == "store ls" {
            getFiles()            
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
