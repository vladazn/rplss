package grpcserver

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"rplss/service/app/pkg/jwt"
)

type AuthInterceptor struct {
	jwt *jwt.JwtPkg
}

func newAuthInterceptor(jwt *jwt.JwtPkg) *AuthInterceptor {
	return &AuthInterceptor{
		jwt: jwt,
	}
}

func (a *AuthInterceptor) authenticate(ctx context.Context, req interface{},
	_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil
	}

	header := md["auth"]
	if len(header) == 0 {
		return handler(ctx, req)
	}

	playerName := a.jwt.ParseKey(header[0])

	ctx = context.WithValue(ctx, "player", playerName)

	return handler(ctx, req)
}
