package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"

    "file_storage_server/server"
    "gorm.io/gorm"
)

// Simple function to ping and test if server is up or not
func getPing(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("pong!\n")
    io.WriteString(w, "pong working!")
}

func postFiles(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    // Limit the size of the uploaded files (optional)
    r.ParseMultipartForm(10 << 20) // 10 MB limit

    // Get all files from the request
    files := r.MultipartForm.File["files"]


    // Loop through the files and save them to the database
    for _, fileHeader := range files {
        // Open the uploaded file
        file, err := fileHeader.Open()
        if err != nil {
            http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        // Read the file content into a byte slice
        fileContent, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error reading file content: %v", err), http.StatusInternalServerError)
            return
        }

        hashDigest := sha256.Sum256(fileContent)
        hashString := hex.EncodeToString(hashDigest[:])

        err = server.CheckDuplicateHash(db, hashString)
        if err != nil {
            http.Error(w, fmt.Sprintf("Content of file %s is already stored in server", fileHeader.Filename), http.StatusBadRequest)
            return
        }

        // Create a new File record to save in the database
        fileInfo := server.File{
            Name:       fileHeader.Filename,
            HashDigest: hashString, 
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
            Content:    string(fileContent), 
        }

        // Save the file record to the database using GORM
        if err := server.CreateFile(db, fileInfo); err != nil {
            http.Error(w, fmt.Sprintf("Error saving file to database: %v", err), http.StatusInternalServerError)
            return
        }
    }

    // Send success response
    fmt.Fprintln(w, "Files uploaded successfully")
}

func main() {
    db, err := server.ConnectToDatabase()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("DB connected\n")


    http.HandleFunc("/ping", getPing)
    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        postFiles(w, r, db) // Pass the db to postFiles
    })

    err = http.ListenAndServe(":2021", nil)

    if err != nil {
        fmt.Printf("Error starting in server: %s\n", err)
        os.Exit(1)
    }

}