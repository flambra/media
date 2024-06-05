package domain

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	Category  string
	Title     string
	Filename  string
	Size      int64
	URL       string
}

type ImageUploadRequest struct {
	UserID   string `form:"user_id"`
	Category string `form:"category"`
	Title    string `form:"title"`
}

type ImageUploadResponse struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
