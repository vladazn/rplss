package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"rplss/service/app/domain"
	"rplss/service/app/pkg/redis"
)

const historyLen = 10

type HistoryRepository struct {
	r redis.Redis
}

func newHistoryRepository(r redis.Redis) *HistoryRepository {
	return &HistoryRepository{r}
}

func (h HistoryRepository) GetHistory(ctx context.Context, player string) []domain.Result {
	response := make([]domain.Result, 0, 10)

	cb := func(item string) error {
		c := domain.Result{}
		err := json.Unmarshal([]byte(item), &c)
		if err != nil {
			return err
		}
		response = append(response, c)
		return nil
	}

	err := h.r.LRange(ctx, h.key(player), cb)
	if err != nil {
		return nil
	}

	return response
}

func (h HistoryRepository) SaveHistory(ctx context.Context, player string, history *domain.Result) {
	key := h.key(player)
	_ = h.r.LPush(ctx, key, history)
	_ = h.r.LTrim(ctx, key, historyLen-1)
}

func (h HistoryRepository) Reset(ctx context.Context, player string) {
	key := h.key(player)
	_ = h.r.Del(ctx, key)
}

func (h HistoryRepository) key(player string) string {
	return fmt.Sprintf("game:%v", player)
}
