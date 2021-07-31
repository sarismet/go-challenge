package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

type RedisDatabase struct {
	Client *redis.Client
}

var (
	Ctx = context.TODO()
)

func NewRedisDatabase() (*RedisDatabase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	if err := client.Ping(Ctx).Err(); err != nil {
		fmt.Println("err ", err)
		return nil, err
	}
	return &RedisDatabase{
		Client: client,
	}, nil
}

func (db *RedisDatabase) InsertKeyToRedis(key string, value string) string {
	return db.Client.Set(Ctx, key, value, 0).Val()
}

func (db *RedisDatabase) GetKeyFromRedis(key string) ([]byte, int) {
	fmt.Println("Hello from Redis")
	val, err := db.Client.Get(Ctx, key).Result()
	if err != nil {
		return nil, 500
	} else if val == "" {
		return nil, 404
	}
	return []byte(val), 0
}

func HelloRedis() {
	fmt.Println("Hello from Redis")
}
