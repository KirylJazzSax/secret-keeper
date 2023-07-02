package main

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/auth/app"
	"github.com/KirylJazzSax/secret-keeper/internal/auth/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/auth"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/user/repository"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	di.ProvideDeps()
	config := do.MustInvoke[*utils.Config](nil)

	switch config.SrvType {
	case commonServer.GRPCType:
		client, err := db.NewMongodbClient(ctx, config)
		if err != nil {
			panic(err)
		}
		defer client.Disconnect(ctx)

		tokenManager := do.MustInvoke[token.Maker](nil)
		repo := repository.NewMongoUserRepository(client)
		hasher := do.MustInvoke[password.PassowrdHasher](nil)

		application := app.NewApplication(
			tokenManager,
			hasher,
			repo,
			config,
		)

		s := server.NewServer(application)

		commonServer.RunGRPCServer(config.GrpcEndpoint, func(srv *grpc.Server) {
			auth.RegisterAuthServiceServer(srv, s)
			reflection.Register(srv)
		})
	case commonServer.GatewayType:
		commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
			auth.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux, config.GrpcEndpoint, opts)
		})
	default:
		panic(errors.UnknownServerType)
	}
}
