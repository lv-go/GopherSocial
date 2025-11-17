package models

type Follower struct {
	FollowedID uint   `json:"followedId"`
	FollowerID uint   `json:"followerId"`
	CreatedAt  string `json:"created_at"`
}
