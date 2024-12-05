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

func pingServer(baseURL string) string {
    serverURL := baseURL + "/ping"

    resp, err := http.Get(serverURL)
    if err != nil {
        log.Println("Error sending request:", err)
        return ""
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error: Received non-OK response status %d\n", resp.StatusCode)
        return ""
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return ""
    }

    fmt.Printf("%s\n", string(body))
    return string(body)
}

func getFiles(baseURL string) string {
    serverURL := baseURL + "/list"

    resp, err := http.Get(serverURL)
    if err != nil {
        log.Println("Error sending request:", err)
        return ""
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error: Received non-OK response status %d\n", resp.StatusCode)
        return ""
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return ""
    }

    fmt.Printf("%s\n", string(body))
    return string(body)
}

func deleteFile(baseURL string, filename string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return "", fmt.Errorf("Error opening %s: %v", filename, err)
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    part, err := writer.CreateFormFile("files", filename)
    if err != nil {
        return "", fmt.Errorf("Error creating form file: %v", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return "", fmt.Errorf("Error copying file content: %v", err)
    }

    err = writer.Close()
    if err != nil {
        return "", fmt.Errorf("Error closing writer: %v", err)
    }

    url := baseURL + "/delete"
    req, err := http.NewRequest("DELETE", url, &requestBody)
    if err != nil {
        return "", fmt.Errorf("Error creating request: %v", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Error: received non-OK response: %v", resp.Status)
    }

    fmt.Println("File successfully deleted!")
    return "File successfully deleted!", nil
}

func putFile(baseURL string, filename string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return "", fmt.Errorf("Error opening %s: %v", filename, err)
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    part, err := writer.CreateFormFile("files", filename)
    if err != nil {
        return "", fmt.Errorf("Error creating form file: %v", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return "", fmt.Errorf("Error copying file content: %v", err)
    }

    err = writer.Close()
    if err != nil {
        return "", fmt.Errorf("Error closing writer: %v", err)
    }

    url := baseURL + "/update"
    req, err := http.NewRequest("PUT", url, &requestBody)
    if err != nil {
        return "", fmt.Errorf("Error creating request: %v", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Error: received non-OK response: %v", resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error updating the file: ", err)
        return "", err
    }

    fmt.Printf("%s\n", string(body))
    return string(body), nil
}

func postFile(baseURL string, filenames []string) (string, error) {
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    for _, filename := range filenames {
        fmt.Println("Attempting to open file:", filename)

        file, err := os.Open(filename)
        if err != nil {
            return "", fmt.Errorf("Error opening file '%s': %v", filename, err)
        }
        defer file.Close()

        part, err := writer.CreateFormFile("files", filename)
        if err != nil {
            return "", fmt.Errorf("Error creating form file for '%s': %v", filename, err)
        }

        _, err = io.Copy(part, file)
        if err != nil {
            return "", fmt.Errorf("Error copying file content for '%s': %v", filename, err)
        }
    }

    err := writer.Close()
    if err != nil {
        return "", fmt.Errorf("Error closing writer: %v", err)
    }

    url := baseURL + "/add"
    req, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        return "", fmt.Errorf("Error creating request: %v", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Error: received non-OK response: %v", resp.Status)
    }

    responseBody, _ := io.ReadAll(resp.Body)
    fmt.Println("Server Response:", string(responseBody))

    fmt.Println("Files created successfully!")
    return "Files created successfully!", nil
}

func getWC(baseURL string) (string, error) {
    serverURL := baseURL + "/wc"

    resp, err := http.Get(serverURL)
    if err != nil {
        log.Println("Error sending request:", err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error: Received non-OK response status %d\n", resp.StatusCode)
        return "", err
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return "", err
    }

    fmt.Printf("%s\n", string(body))
    return string(body), nil
}

func getFW(baseURL string, limit string, order string) (string, error) {
    serverURL := baseURL + "/fw?limit=" + limit + "&order=" + order

    resp, err := http.Get(serverURL)
    if err != nil {
        log.Println("Error sending request:", err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error: Received non-OK response status %d\n", resp.StatusCode)
        return "", err
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return "", err
    }

    fmt.Printf("%s\n", string(body))
    return string(body), nil
}


func main() {
    baseURL := "http://localhost:2021"

    fmt.Println("CLI Program started. Type 'store' to send a request to the server.")

    scanner := bufio.NewScanner(os.Stdin)

    for {
        // Print to show newline in which user can put command
        fmt.Print("> ")
        scanner.Scan()

        command := strings.TrimSpace(scanner.Text())

        if strings.HasPrefix(command, "store rm ") {
            filename := strings.TrimPrefix(command, "store rm ")
            fmt.Printf("Sending delete request for file: %s\n", filename)
            _, err := deleteFile(baseURL, filename)
            if err != nil {
                log.Printf("Error: %v\n", err)
            }
        } else if strings.HasPrefix(command, "store update ") {
            filename := strings.TrimPrefix(command, "store update ")
            fmt.Printf("Sending update request for file: %s\n", filename)
            _, err := putFile(baseURL, filename)
            if err != nil {
                log.Printf("Error: %v\n", err)
            }
        } else if strings.HasPrefix(command, "store add") {
            parts := strings.Fields(command)
            filenames := parts[2:]
            fmt.Println("Sending create request")
            _, err := postFile(baseURL, filenames)
            if err != nil {
                log.Printf("Error: %v\n", err)
            }
        } else if command == "store wc" {
            getWC(baseURL)
        } else if command == "store" {
            fmt.Println("Sending request to the server...")
            pingServer(baseURL)
        } else if command == "store ls" {
            getFiles(baseURL) 
        } else if strings.HasPrefix(command, "store freq-words") {
            parts := strings.Fields(command)
            limit := parts[3]
            order := strings.Split(parts[4], "=")[1]
            getFW(baseURL, limit, order)
        } else if command == "exit" {
            fmt.Println("Exiting program...")
            break
        } else {
            fmt.Println("Unknown command:", command)
        }
    }
}
