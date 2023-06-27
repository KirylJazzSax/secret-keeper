package main

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/auth/application"
	"github.com/KirylJazzSax/secret-keeper/internal/auth/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/auth"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/user/repository"
	"github.com/boltdb/bolt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := do.MustInvoke[*utils.Config](nil)
	b, err := bolt.Open(config.DbUrl, 0600, nil)
	tokenManager := do.MustInvoke[token.Maker](nil)
	repo := repository.NewUserRepository(b)
	hasher := do.MustInvoke[password.PassowrdHasher](nil)

	app := application.NewApplication(
		tokenManager,
		hasher,
		repo,
		config,
	)

	s := server.NewServer(app)

	endpoint := fmt.Sprintf(":%s", config.Port)

	go commonServer.RunGRPCServer(endpoint, func(srv *grpc.Server) {
		auth.RegisterAuthServiceServer(srv, s)
		reflection.Register(srv)
	})

	commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
		auth.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
	})

}

func setupDb(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.UsersBucket))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
}
