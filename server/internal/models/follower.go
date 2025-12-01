package models

import (
	"time"
)

type Follower struct {
	FollowedID string    `json:"followedId" gorm:"type:uuid"`
	FollowerID string    `json:"followerId" gorm:"type:uuid"`
	CreatedAt  time.Time `json:"created_at"`
}
