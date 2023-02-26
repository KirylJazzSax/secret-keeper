package main

import (
	"fmt"
	"net"
	"secret_keeper/db"
	"secret_keeper/errors"
	"secret_keeper/gapi"
	"secret_keeper/internal/di"
	"secret_keeper/pb"
	"secret_keeper/utils"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := di.ProvideDeps("."); err != nil {
		panic(err)
	}
	config := do.MustInvoke[*utils.Config](nil)

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.PORT))

	if err != nil {
		errors.LogErr(err)
	}

	boltDb, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		errors.LogErr(err)
	}
	defer boltDb.Close()

	err = db.SetupDb(boltDb)
	if err != nil {
		errors.LogErr(err)
	}

	s, err := gapi.NewServer(db.NewBoltStore(boltDb), config)
	if err != nil {
		errors.LogErr(err)
	}

	pb.RegisterSecretKeeperServer(grpcServer, s)
	reflection.Register(grpcServer)

	log.Info().Msg("server runs")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Stack().Err(err)
	}
}
