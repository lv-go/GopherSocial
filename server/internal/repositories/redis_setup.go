package repositories

import "github.com/go-redis/redis/v8"

var redisClient *redis.Client

type RedisConfig struct {
	Addr    string
	Pw      string
	Db      int
	Enabled bool
}

func SetupRedisClient(cfg RedisConfig) *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pw,
		DB:       cfg.Db,
	})
	return redisClient
}
