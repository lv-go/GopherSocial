package models

type Follower struct {
	ID        int64  `json:"id"`
	UserID    uint   `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
