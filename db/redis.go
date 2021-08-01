package db

import (
	"context"
	"os"

	"github.com/go-redis/redis"
)

type RedisDatabase struct {
	Client *redis.Client
}

var (
	Ctx = context.TODO()
)

//Creates a new redis client
func NewRedisDatabase() (*RedisDatabase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisDatabase{
		Client: client,
	}, nil
}

//Insert a key value pair to redis
func (db *RedisDatabase) InsertKeyToRedis(key string, value string) string {
	return db.Client.Set(Ctx, key, value, 0).Val()
}

//Gets the value using its key from redis
func (db *RedisDatabase) GetKeyFromRedis(key string) ([]byte, int) {
	val, err := db.Client.Get(Ctx, key).Result()
	if val == "" || (err != nil && err.Error() == "redis: nil") {
		return nil, 404
	} else if err != nil {
		return nil, 500
	}
	return []byte(val), 0
}
