package repositories

import "github.com/go-redis/redis/v8"

var redisDB *redis.Client

type RedisConfig struct {
	Addr    string
	Pw      string
	Db      int
	Enabled bool
}

func SetupRedisDB(cfg RedisConfig) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pw,
		DB:       cfg.Db,
	})
}
