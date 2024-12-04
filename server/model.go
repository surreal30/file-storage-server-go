package server

import (
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
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