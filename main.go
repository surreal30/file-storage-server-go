package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
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


func getFiles(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    // Fetch all files from the database
    files, err := server.GetFiles(db)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching files: %v", err), http.StatusInternalServerError)
        return
    }

    for _, file := range files {
        fmt.Fprintf(w, "File ID: %d, Name: %s \n", file.ID, file.Name)
    }
}

func deleteFile(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    r.ParseMultipartForm(10 << 20)
    files := r.MultipartForm.File["files"]

    for _, fileHeader := range files {

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

        err = server.DeleteFile(db, hashString)
        if err != nil {
            http.Error(w, fmt.Sprintf("Some error occured while deleting, %s", err), http.StatusInternalServerError)
            return
        }
    }

    fmt.Fprintln(w, "File deleted successfully")
}

func putFile(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    r.ParseMultipartForm(10 << 20)
    files := r.MultipartForm.File["files"]

    for _, fileHeader := range files {

        uploadedFile, err := fileHeader.Open()
        if err != nil {
            http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
            return
        }
        defer uploadedFile.Close()

        // Read the file content into a byte slice
        fileContent, err := io.ReadAll(uploadedFile)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error reading file content: %v", err), http.StatusInternalServerError)
            return
        }

        hashDigest := sha256.Sum256(fileContent)
        hashString := hex.EncodeToString(hashDigest[:])

        existingFile, err := server.GetFileByName(db, fileHeader.Filename)

        // If err is nil then the file has been found 
        if err == nil {
            existingFile.Content = string(fileContent)
            existingFile.HashDigest = hashString
            existingFile.UpdatedAt = time.Now()

            err = server.UpdateFile(db, existingFile)
            if err != nil {
                http.Error(w, fmt.Sprintf("Some error occured while updating, %s", err), http.StatusInternalServerError)
                return
            }
            fmt.Fprintln(w, "Files updated successfully")
            return
        }

        // If file is not found
        fileInfo := server.File{
            Name:       fileHeader.Filename,
            HashDigest: hashString, 
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
            Content:    string(fileContent), 
        }

        if err := server.CreateFile(db, fileInfo); err != nil {
            http.Error(w, fmt.Sprintf("Error saving file to database: %v", err), http.StatusInternalServerError)
            return
        }
    }

    fmt.Fprintln(w, "Files uploaded successfully")
}

func getWordCount(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    content, err := server.FetchContentAllFile(db)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching files: %v", err), http.StatusInternalServerError)
        return
    }

    wc := len(strings.Split(content, " "))
    fmt.Fprintf(w, "All files contain %d words \n", wc)
}

func main() {
    db, err := server.ConnectToDatabase()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("DB connected\n")

    http.HandleFunc("/ping", getPing)
    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        postFiles(w, r, db)
    })
    http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        getFiles(w, r, db)
    })
    http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
        deleteFile(w, r, db)
    })
    http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
        putFile(w, r, db)
    })
    http.HandleFunc("/wc", func(w http.ResponseWriter, r *http.Request) {
        getWordCount(w, r, db)
    })


    err = http.ListenAndServe(":2021", nil)

    if err != nil {
        fmt.Printf("Error starting in server: %s\n", err)
        os.Exit(1)
    }

}