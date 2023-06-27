package main

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/user"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app"
	"github.com/KirylJazzSax/secret-keeper/internal/user/repository"
	"github.com/KirylJazzSax/secret-keeper/internal/user/server"

	"github.com/boltdb/bolt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	di.ProvideDeps(".")
	config := do.MustInvoke[*utils.Config](nil)

	b, err := bolt.Open(config.DbUrl, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer b.Close()

	if err = setupDb(b); err != nil {
		panic(err)
	}

	hasher := do.MustInvoke[password.PassowrdHasher](nil)
	validator := do.MustInvoke[validation.Validator](nil)
	repo := repository.NewUserRepository(b)

	a := app.NewApplication(validator, hasher, repo)
	s := server.NewServer(a)

	endpoint := fmt.Sprintf(":%s", config.Port)

	go commonServer.RunGRPCServer(endpoint, func(srv *grpc.Server) {
		user.RegisterUsersServiceServer(srv, s)
		reflection.Register(srv)
	})

	commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
		user.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
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
