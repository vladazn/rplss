package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"rplss/service/app/domain"
	"rplss/service/app/pkg/redis"
	"rplss/service/config"
	"testing"
)

const configPath = "../../config/config.yml"

func shouldTest() bool {
	return os.Getenv("MANUAL") != ""
}

func TestHistory(t *testing.T) {
	if !shouldTest() {
		t.Skip()
		return
	}

	ctx := context.Background()
	_ = ctx

	var err error
	configs, err := config.New(configPath)
	require.NoError(t, err)

	r, err := redis.NewRedisConnection(configs.Redis)
	require.NoError(t, err)

	h := newHistoryRepository(r)

	player1 := "player1"
	player2 := "player2"

	history1 := domain.Result{
		Results:  "win",
		Player:   1,
		Computer: 2,
	}

	history2 := domain.Result{
		Results:  "lose",
		Player:   1,
		Computer: 2,
	}

	h.SaveHistory(ctx, player1, &history1)
	h.SaveHistory(ctx, player1, &history2)
	h.SaveHistory(ctx, player2, &history2)

	histories := h.GetHistory(ctx, player1)
	require.Equal(t, 2, len(histories))
	require.Equal(t, "lose", histories[0].Results)

}
