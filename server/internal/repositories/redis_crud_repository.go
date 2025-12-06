package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCRUDRepository[T interface{}, ID any] struct {
	redisClient  *redis.Client
	cacheBaseKey string
	expTime      time.Duration
}

func NewRedisCRUDRepository[T interface{}, ID any](
	cacheBaseKey string,
	expTime time.Duration,
) CRUDRepository[T, ID] {
	return &redisCRUDRepository[T, ID]{
		redisClient:  redisClient,
		cacheBaseKey: cacheBaseKey,
		expTime:      expTime,
	}
}

func (r *redisCRUDRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	cacheKey := fmt.Sprintf("%s-%v", r.cacheBaseKey, entity)
	jsonEntity, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	return r.redisClient.SetEX(ctx, cacheKey, jsonEntity, r.expTime).Err()
}

func (r *redisCRUDRepository[T, ID]) GetByID(ctx context.Context, id ID) (*T, error) {
	cacheKey := fmt.Sprintf("%s-%v", r.cacheBaseKey, id)

	data, err := r.redisClient.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var entity T
	if data != "" {
		err := json.Unmarshal([]byte(data), &entity)
		if err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *redisCRUDRepository[T, ID]) GetOne(ctx context.Context, filter interface{}) (*T, error) {
	cacheKey := fmt.Sprintf("%s-%v", r.cacheBaseKey, filter)

	data, err := r.redisClient.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var entity T
	if data != "" {
		err := json.Unmarshal([]byte(data), &entity)
		if err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *redisCRUDRepository[T, ID]) GetAll(ctx context.Context, filter interface{}) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCRUDRepository[T, ID]) GetPage(ctx context.Context, filter interface{}, pageQuery PageQuery) (*Page[T], error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCRUDRepository[T, ID]) UpdateByID(ctx context.Context, id ID, entity *T) error {
	cacheKey := fmt.Sprintf("%s-%v", r.cacheBaseKey, id)

	jsonEntity, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	return r.redisClient.SetEX(ctx, cacheKey, jsonEntity, r.expTime).Err()
}

func (r *redisCRUDRepository[T, ID]) DeleteByID(ctx context.Context, id ID) error {
	cacheKey := fmt.Sprintf("%s-%v", r.cacheBaseKey, id)
	intCmd := r.redisClient.Del(ctx, cacheKey)
	if intCmd.Err() != nil {
		return intCmd.Err()
	}
	return nil
}

func (r *redisCRUDRepository[T, ID]) DeleteOne(ctx context.Context, filter interface{}) error {
	//TODO implement me
	panic("implement me")
}
