package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Password       string         `json:"-"`
	IsActive       bool           `json:"is_active"`
	RoleID         int64          `json:"role_id"`
	Role           Role           `json:"role"`
	Posts          []Post         `json:"posts"`
	Comments       []Comment      `json:"comments" gorm:"foreignKey:UserID"`
	UserInvitation UserInvitation `json:"user_invitation" gorm:"foreignKey:UserID"`
}

type UserInvitation struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}
