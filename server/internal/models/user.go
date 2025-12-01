package models

import (
	"time"
)

type User struct {
	ID             string         `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Password       string         `json:"-"`
	IsActive       bool           `json:"is_active"`
	RoleID         int64          `json:"role_id"`
	Role           *Role          `json:"role"`
	Posts          []Post         `json:"posts"`
	Comments       []Comment      `json:"comments" gorm:"foreignKey:UserID"`
	UserInvitation UserInvitation `json:"user_invitation" gorm:"foreignKey:UserID"`
}

type UserInvitation struct {
	UserID string    `json:"userId" gorm:"type:uuid"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
