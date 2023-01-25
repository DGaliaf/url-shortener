package reduction

import (
	"context"
	"url-shorter/app/internal/adapters/db/redis"
	"url-shorter/app/internal/domain/entities"
)

type Storage interface {
	Add(ctx context.Context, reduction entities.Reduction) (string, error)
}

type ReductionService struct {
	storage *redis.ReductionStorage
}

func NewReductionService(storage *redis.ReductionStorage) *ReductionService {
	return &ReductionService{storage: storage}
}

func (r ReductionService) CreateShortUrl(ctx context.Context, dto CreateShortUrlDTO) (string, error) {
	reduction := entities.Reduction{
		LongUrl:    dto.LongUrl,
		CustomText: dto.CustomText,
	}

	if err := reduction.Validate(); err != nil {
		return "", err
	}

	return r.storage.Add(ctx, reduction)
}

func (r ReductionService) GetLongUrl(ctx context.Context, hashedLink string) (string, error) {
	return r.storage.Get(ctx, hashedLink)
}
