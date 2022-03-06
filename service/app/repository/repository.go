package repository

import (
	"context"
	"rplss/service/app/domain"
	"rplss/service/app/pkg/redis"
)

type History interface {
	GetHistory(ctx context.Context, player string) []domain.Result
	SaveHistory(ctx context.Context, player string, history *domain.Result)
	Reset(ctx context.Context, player string)
}

type Repositories struct {
	History History
}

func InitRepositories(r redis.Redis) *Repositories {
	return &Repositories{History: newHistoryRepository(r)}
}
