package domain

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	URL        string         `gorm:"type:text;not null" json:"url"`
	UploadedAt time.Time      `gorm:"type:timestamp;not null;default:current_timestamp" json:"uploaded_at"`
}

type ImageUploadRequest struct {
	Title string `form:"title"`
}
