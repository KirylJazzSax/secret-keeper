package main

import (
	"fmt"
	"net"
	"secret_keeper/errors"
	"secret_keeper/gapi"
	"secret_keeper/internal/di"
	"secret_keeper/pb"
	"secret_keeper/repository"
	"secret_keeper/utils"

	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := di.ProvideDeps("."); err != nil {
		panic(err)
	}
	defer do.Shutdown[repository.Repository](nil)

	server := do.MustInvoke[*gapi.Server](nil)

	config := do.MustInvoke[*utils.Config](nil)

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.PORT))
	if err != nil {
		errors.LogErr(err)
	}

	pb.RegisterSecretKeeperServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Info().Msg("server runs")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Stack().Err(err)
	}
}
