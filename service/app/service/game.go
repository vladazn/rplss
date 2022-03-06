package service

import (
	"context"
	"math/rand"
	"rplss/service/app/domain"
	"rplss/service/app/pkg/jwt"
	"rplss/service/app/repository"
)

type GameService struct {
	repo    *repository.Repositories
	rules   map[string][]string
	choices []string
	jwt     *jwt.JwtPkg
}

func newGameService(repo *repository.Repositories, rules map[string][]string,
	jwt *jwt.JwtPkg) *GameService {
	choices := make([]string, 0, len(rules))
	for choice := range rules {
		choices = append(choices, choice)
	}

	return &GameService{
		repo:    repo,
		choices: choices,
		rules:   rules,
		jwt:     jwt,
	}
}

func (g GameService) Choices() []string {
	return g.choices
}

func (g GameService) Choice() *domain.Choice {
	n := rand.Intn(len(g.choices))

	return &domain.Choice{
		Id:   n + 1,
		Name: g.choices[n],
	}
}

func (g GameService) Play(ctx context.Context, choice int, player string) *domain.Result {
	n := rand.Intn(len(g.choices))

	result := &domain.Result{
		Player:   choice,
		Computer: n + 1,
	}

	choice = choice - 1

	switch true {
	case n == choice:
		result.Results = "tie"
	case g.isWin(choice, n):
		result.Results = "win"
	default:
		result.Results = "lose"
	}

	if player != "" {
		g.repo.History.SaveHistory(ctx, player, result)
	}
	g.repo.History.SaveHistory(ctx, "global", result)

	return result
}

func (g GameService) GetHistory(ctx context.Context, player string) []domain.Result {

	if player != "" {
		return g.repo.History.GetHistory(ctx, player)
	} else {
		return g.repo.History.GetHistory(ctx, "global")
	}
}

func (g GameService) ResetHistory(ctx context.Context, player string) {
	g.repo.History.Reset(ctx, player)
}

func (g GameService) Login(player string) string {
	return g.jwt.NewKey(player)
}

func (g GameService) isWin(p1, p2 int) bool {
	p2Name := g.choices[p2]

	for _, p1Name := range g.rules[g.choices[p1]] {
		if p1Name == p2Name {
			return true
		}
	}

	return false
}
