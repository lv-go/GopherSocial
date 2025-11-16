package repositories

import "github.com/go-redis/redis/v8"

var redisDB *redis.Client

type RedisConfig struct {
	Addr    string
	Pw      string
	Db      int
	Enabled bool
}

func SetupRedisDB(addr, pw string, db int) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}
