package main

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/user"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app"
	"github.com/KirylJazzSax/secret-keeper/internal/user/repository"
	"github.com/KirylJazzSax/secret-keeper/internal/user/server"
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

		hasher := do.MustInvoke[password.PassowrdHasher](nil)
		validator := do.MustInvoke[validation.Validator](nil)
		repo := repository.NewMongoUserRepository(client)

		a := app.NewApplication(validator, hasher, repo)
		s := server.NewServer(a)
		commonServer.RunGRPCServer(config.GrpcEndpoint, func(srv *grpc.Server) {
			user.RegisterUsersServiceServer(srv, s)
			reflection.Register(srv)
		})
	case commonServer.GatewayType:
		fmt.Println("gateway")
		commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
			user.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, opts)
		})
	default:
		panic(errors.UnknownServerType)
	}
}
