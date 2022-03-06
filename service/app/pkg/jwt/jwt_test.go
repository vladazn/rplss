package jwt

import (
	"context"
	"github.com/stretchr/testify/require"
	"rplss/service/config"
	"testing"
)

const configPath = "../../../config/config.yml"

func TestJwt(t *testing.T) {
	ctx := context.Background()
	_ = ctx

	var err error
	configs, err := config.New(configPath)
	require.NoError(t, err)

	j := NewJwtPkg(configs.JWT.Key)

	playerName := "aaa"

	key := j.NewKey(playerName)
	require.NotEmpty(t, key)

	player := j.ParseKey(key)
	require.Equal(t, playerName, player)
}
