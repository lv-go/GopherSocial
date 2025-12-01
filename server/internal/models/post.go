package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model `json:"-"`
	ID         uint           `json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Content    string         `json:"content"`
	Title      string         `json:"title"`
	UserID     string         `json:"userId" gorm:"type:uuid"`
	Tags       pq.StringArray `json:"tags" gorm:"type:text[]"`
	Version    int            `json:"version"`
	Comments   []Comment      `json:"comments" gorm:"foreignKey:PostID"`
}
