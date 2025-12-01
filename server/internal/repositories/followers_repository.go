package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type FollowersRepository struct {
	db *gorm.DB
}

func NewFollowersRepository() *FollowersRepository {
	return &FollowersRepository{
		db: gormDB,
	}
}

func (r *FollowersRepository) Follow(ctx context.Context, followedID, followerID string) error {
	return r.db.WithContext(ctx).Create(&models.Follower{
		FollowedID: followedID,
		FollowerID: followerID,
	}).Error
}

func (r *FollowersRepository) Unfollow(ctx context.Context, followedID, followerID string) error {
	return r.db.WithContext(ctx).Delete(&models.Follower{
		FollowedID: followedID,
		FollowerID: followerID,
	}).Error
}
