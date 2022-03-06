package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	grpcserver "rplss/service/app/pkg/gprc/server"
	"rplss/service/app/pkg/jwt"
	"rplss/service/app/pkg/redis"
	"rplss/service/app/repository"
	"rplss/service/app/service"
	"rplss/service/config"
	"syscall"
)

func Run(configPath string) {
	ctx := context.Background()
	_ = ctx

	var err error
	configs, err := config.New(configPath)
	if err != nil {
		return
	}

	r, err := redis.NewRedisConnection(configs.Redis)
	if err != nil {
		return
	}

	rp := repository.InitRepositories(r)

	j := jwt.NewJwtPkg(configs.JWT.Key)

	services := service.InitServices(rp, configs.Choices, j)

	server := grpcserver.NewGrpcServer(configs.Grpc.Host, configs.Grpc.Port, services, j)

	errCh := make(chan error)
	go func() {
		err := server.Serve()

		if err != nil {
			errCh <- err
		}
	}()

	fmt.Printf("grpc server started on port: %v\n", configs.Grpc.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		fmt.Println(err)
	case <-quit:
		fmt.Println("quit call")
	}

	fmt.Println("stopping")

	_ = r.Close()
	server.Stop()

}
