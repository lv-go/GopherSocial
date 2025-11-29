package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"userId" gorm:"type:uuid"`
	User    User      `json:"user"`
	PostID  uint      `json:"post_id"`
	Post    Post      `json:"post"`
}
