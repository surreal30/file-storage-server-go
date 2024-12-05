package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "sort"
    "strconv"
    "strings"
    "time"

    "file_storage_server/server"
    "gorm.io/gorm"
)

type WordCount struct {
    Word  string `json:"word"`
    Count int    `json:"count"`
}


// Simple function to ping and test if server is up or not
func getPing(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("pong!\n")
    io.WriteString(w, "pong working!")
}

// Save files in DB
func postFiles(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    r.ParseMultipartForm(10 << 20) 

    files := r.MultipartForm.File["files"]

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
            return
        }
        defer file.Close()

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

    fmt.Fprintln(w, "Files uploaded successfully")
}

// Get list of files
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

// Delete file
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

// Update a file if it exists otherwise create a new file
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

// Fetch word count
func getWordCount(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    content, err := server.FetchContentAllFile(db)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching files: %v", err), http.StatusInternalServerError)
        return
    }

    wc := len(strings.Split(content, " "))
    fmt.Fprintf(w, "All files contain %d words \n", wc)
}

// Fetch frequent words
func getFreqWord(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    limitStr := r.URL.Query().Get("limit")
    order := r.URL.Query().Get("order")

    limit := 5

    fmt.Println(limitStr, order)

    if limitStr != "" {
        fmt.Println(limitStr, order)

        parsedLimit, err := strconv.Atoi(limitStr)
        fmt.Println("helolsosbfk", err)

        if err != nil {
            http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
            return
        }
        limit = parsedLimit
    }

    fmt.Println("ordering now")

    if order != "asc" && order != "dsc" {
        http.Error(w, "Invalid 'order' parameter. Use 'asc' or 'dsc'.", http.StatusBadRequest)
        return
    }

    content, err := server.FetchContentAllFile(db)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching files: %v", err), http.StatusInternalServerError)
        return
    }

    fmt.Println("files fetched now")


    words := strings.Split(content, " ")
    wordCounts := make(map[string]int)

    // Count frequency of each word
    for _, word := range words {
        wordCounts[word]++
    }

    var wordCountList []WordCount
    for word, count := range wordCounts {
        wordCountList = append(wordCountList, WordCount{
            Word:  word,
            Count: count,
        })
    }

    fmt.Println("counting now")

    var ordering string
    if order == "dsc" {
        sort.Slice(wordCountList, func(i, j int) bool {
            return wordCountList[i].Count > wordCountList[j].Count
        })
        ordering = "most"
    } else {
        sort.Slice(wordCountList, func(i, j int) bool {
            return wordCountList[i].Count < wordCountList[j].Count
        })
        ordering = "least"
    }

    fmt.Println("sorting now")


    if limit > len(wordCountList) {
        limit = len(wordCountList)
    }

    fmt.Fprintf(w, "The %d %s frequent words are:\n", limit, ordering)
    for _, wc := range wordCountList[:limit] {
        fmt.Fprintf(w, "%s\n", wc.Word)
    }
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
    http.HandleFunc("/fw", func(w http.ResponseWriter, r *http.Request) {
        getFreqWord(w, r, db)
    })


    err = http.ListenAndServe(":2021", nil)

    if err != nil {
        fmt.Printf("Error starting in server: %s\n", err)
        os.Exit(1)
    }

}