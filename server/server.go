package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

type File struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Path       string    `gorm:"type:varchar(255);not null"`
	HashDigest string    `gorm:"type:varchar(256)"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"type:datetime"`
	UpdatedAt  time.Time `gorm:"type:datetime"`
}

// Simple function to ping and test if server is up or not
func getPing(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("pong!\n")
	io.WriteString(w, "pong working!")
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Connected to the database successfully!")


	http.HandleFunc("/ping", getPing)

	err = http.ListenAndServe(":2021", nil)

	if err != nil {
		fmt.Printf("Error starting in server: %s\n", err)
		os.Exit(1)
	}

}