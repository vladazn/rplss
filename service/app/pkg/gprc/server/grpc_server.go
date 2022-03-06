package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"rplss/proto/gen/go/proto/game"
	"rplss/service/app/pkg/jwt"
	"rplss/service/app/service"
)

type Server struct {
	host     string
	port     int
	server   *grpc.Server
	services *service.Services
}

func NewGrpcServer(
	host string,
	port int,
	services *service.Services,
	jwt *jwt.JwtPkg,
) Server {

	auth := newAuthInterceptor(jwt)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(auth.authenticate),
	}

	return Server{
		host:     host,
		port:     port,
		server:   grpc.NewServer(opts...),
		services: services,
	}
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}

	game.RegisterGameServer(s.server, newGameServer(s.services))

	err = s.server.Serve(lis)

	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
