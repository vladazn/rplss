package grpcserver

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"rplss/proto/gen/go/proto/game"
	"rplss/service/app/service"
)

type GameServer struct {
	game.UnimplementedGameServer
	services *service.Services
}

func (g *GameServer) Choices(_ context.Context, _ *emptypb.Empty) (*game.ChoicesResponse,
	error) {
	choices := g.services.Game.Choices()
	response := &game.ChoicesResponse{
		Choices: make([]*game.ChoiceResponse, len(choices)),
	}

	for i, c := range choices {
		response.Choices[i] = &game.ChoiceResponse{
			Id:   int32(i + 1),
			Name: c,
		}
	}

	return response, nil
}

func (g *GameServer) Choice(_ context.Context, _ *emptypb.Empty) (*game.ChoiceResponse,
	error) {
	c := g.services.Game.Choice()

	return &game.ChoiceResponse{
		Id:   int32(c.Id),
		Name: c.Name,
	}, nil
}

func (g *GameServer) Login(_ context.Context, req *game.LoginRequest) (*game.LoginResponse,
	error) {
	c := g.services.Game.Login(req.Username)

	return &game.LoginResponse{
		Jwt: c,
	}, nil
}

func (g *GameServer) Play(ctx context.Context, req *game.PlayRequest) (*game.PlayResponse,
	error) {

	v := ctx.Value("player")
	playerName := ""
	if v != nil {
		playerName = v.(string)
	}

	res := g.services.Game.Play(ctx, int(req.Player), playerName)

	return &game.PlayResponse{
		Results:  res.Results,
		Player:   int32(res.Player),
		Computer: int32(res.Computer),
	}, nil
}

func (g *GameServer) Reset(ctx context.Context, _ *emptypb.Empty) (*game.SuccessResponse, error) {
	v := ctx.Value("player")
	playerName := ""
	if v != nil {
		playerName = v.(string)
	}

	if playerName == "" {
		return &game.SuccessResponse{
			Success: false,
		}, nil
	}

	g.services.Game.ResetHistory(ctx, playerName)

	return &game.SuccessResponse{
		Success: true,
	}, nil
}

func (g *GameServer) History(ctx context.Context, _ *emptypb.Empty) (*game.HistoryResponse, error) {
	v := ctx.Value("player")
	playerName := ""
	if v != nil {
		playerName = v.(string)
	}

	histories := g.services.Game.GetHistory(ctx, playerName)

	r := &game.HistoryResponse{
		Results: make([]*game.PlayResponse, len(histories)),
	}

	for i, h := range histories {
		r.Results[i] = &game.PlayResponse{
			Results:  h.Results,
			Player:   int32(h.Player),
			Computer: int32(h.Computer),
		}
	}

	return r, nil
}

func newGameServer(services *service.Services) *GameServer {
	return &GameServer{services: services}
}
