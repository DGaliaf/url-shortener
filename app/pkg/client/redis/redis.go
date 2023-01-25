package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
)

type redisConfig struct {
	Host string
	Port string
	DB   int
}

func NewRedisConfig(host string, port string, DB int) *redisConfig {
	return &redisConfig{Host: host, Port: port, DB: DB}
}

func NewClient(ctx context.Context, cfg *redisConfig) *redis.Client {
	log.Println("client initializing")
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DB:   cfg.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	return rdb
}
