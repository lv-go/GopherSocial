package models

import (
	"time"

	"github.com/google/uuid"
)

type Follower struct {
	FollowedID uuid.UUID `json:"followedId" gorm:"type:uuid"`
	FollowerID uuid.UUID `json:"followerId" gorm:"type:uuid"`
	CreatedAt  time.Time `json:"created_at"`
}
