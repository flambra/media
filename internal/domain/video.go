package domain

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UserID      uint
	Category    string
	Title       string
	Description string
	Filename    string
	Size        int64
	S3Url       string
	URL         string
}

type VideoUploadRequest struct {
	UserID      uint   `form:"user_id"`
	Category    string `form:"category"`
	Title       string `form:"title"`
	Description string `form:"description"`
}

type VideoUploadResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
