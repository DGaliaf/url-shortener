package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"url-shorter/app/internal/config"
	"url-shorter/app/internal/domain/entities"
	"url-shorter/app/pkg/shortener"
)

type ReductionStorage struct {
	cfg *config.Config
	db  *redis.Client
}

func NewReductionStorage(cfg *config.Config, db *redis.Client) *ReductionStorage {
	return &ReductionStorage{cfg: cfg, db: db}
}

func (r ReductionStorage) Add(ctx context.Context, reduction entities.Reduction) (string, error) {
	link := shortener.GenerateShortLink(reduction.LongUrl)

	cachedUrl, err := r.Get(ctx, link)
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
	}

	shortUrl := fmt.Sprintf("http://%s:%d/%s", r.cfg.HTTP.IP, r.cfg.HTTP.Port, link)

	if r.cfg.HTTP.IP == "0.0.0.0" {
		shortUrl = fmt.Sprintf("http://localhost:%d/%s", r.cfg.HTTP.Port, link)
	}

	if cachedUrl != "" {
		return shortUrl, nil
	}

	if err := r.db.Set(ctx, link, reduction.LongUrl, 0).Err(); err != nil {
		return "", err
	}

	return shortUrl, nil
}

func (r ReductionStorage) Get(ctx context.Context, link string) (string, error) {
	result, err := r.db.Get(ctx, link).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
