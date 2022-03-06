package service

import (
	"context"
	"rplss/service/app/domain"
	"rplss/service/app/pkg/jwt"
	"rplss/service/app/repository"
)

type Game interface {
	Choices() []string
	Choice() *domain.Choice
	Play(ctx context.Context, choice int, player string) *domain.Result
	GetHistory(ctx context.Context, player string) []domain.Result
	ResetHistory(ctx context.Context, player string)
	Login(nickname string) string
}

type Services struct {
	Game Game
}

func InitServices(repo *repository.Repositories, rules map[string][]string,
	jwt *jwt.JwtPkg) *Services {

	return &Services{
		Game: newGameService(repo, rules, jwt),
	}
}
