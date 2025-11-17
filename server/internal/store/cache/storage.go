package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sikozonpc/social/internal/models"
)

type Storage struct {
	Users interface {
		Get(context.Context, uint) (*models.User, error)
		Set(context.Context, *models.User) error
		Delete(context.Context, uint)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}
