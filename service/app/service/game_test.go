package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"rplss/service/app/pkg/jwt"
	"rplss/service/app/pkg/redis"
	"rplss/service/app/repository"
	"rplss/service/config"
	"testing"
)

const configPath = "../../config/config.yml"

func shouldTest() bool {
	return os.Getenv("MANUAL") != ""
}

func TestGame(t *testing.T) {
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

	rp := repository.InitRepositories(r)

	j := jwt.NewJwtPkg(configs.JWT.Key)

	g := newGameService(rp, configs.Choices, j)

	c := g.Choices()
	require.Equal(t, 5, len(c))

	single := g.Choice()
	require.NotNil(t, single)

	player := "player1"
	//g.ResetHistory(ctx, player)

	g.Play(ctx, 2, player)
	g.Play(ctx, 2, player)

	histories := g.GetHistory(ctx, player)
	require.Equal(t, 2, len(histories))

	g.ResetHistory(ctx, player)

	require.True(t, g.isWin(0, 2))
}
