package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type FollowersRepository struct {
	db *gorm.DB
}

func NewFollowersRepository(db *gorm.DB) *FollowersRepository {
	return &FollowersRepository{
		db: db,
	}
}

func (r *FollowersRepository) Follow(ctx context.Context, followedID uint, followerID uint) error {
	return r.db.WithContext(ctx).Create(&models.Follower{
		FollowedID: followedID,
		FollowerID: followerID,
	}).Error
}

func (r *FollowersRepository) Unfollow(ctx context.Context, followedID uint, followerID uint) error {
	return r.db.WithContext(ctx).Delete(&models.Follower{
		FollowedID: followedID,
		FollowerID: followerID,
	}).Error
}
