package server

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)



func ConnectToDatabase() (*gorm.DB, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
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
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	return db, nil
}

func CreateFile(db *gorm.DB, file File) error {
	fmt.Printf(file.Name)
	// Create the new file record
	if err := db.Create(&file).Error; err != nil {
		return err
	}
	return nil
}

func GetFiles(db *gorm.DB) ([]File, error) {
	var files []File

	result := db.Find(&files)  // Retrieve all records from the 'files' table
    if result.Error != nil {
        return nil, result.Error  // Return the error if something goes wrong
    }
    return files, nil 
}

func DeleteFile(db *gorm.DB, key string) (error) {
	var file File

	result := db.Where("hash_digest = ?", key).Delete(&file)
	if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("file not found")
    }

    return nil
}

func CheckDuplicateHash(db *gorm.DB, hashDigest string) error {
    var file File
    // Query to find if a file with the given hash_digest exists
    result := db.Where("hash_digest = ?", hashDigest).First(&file)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            // If no record is found, it's not a duplicate, return nil
            return nil
        }
        // If there is any other error, return it
        return result.Error
    }
    // If the file is found, it's a duplicate, return an error
    return fmt.Errorf("file with hash_digest %s already exists", hashDigest)
}