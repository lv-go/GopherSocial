package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Content  string         `json:"content"`
	Title    string         `json:"title"`
	UserID   uint           `json:"user_id"`
	Tags     pq.StringArray `json:"tags" gorm:"type:text[]"`
	Version  int            `json:"version"`
	Comments []Comment      `json:"comments" gorm:"foreignKey:PostID"`
	User     User           `json:"user"`
}
