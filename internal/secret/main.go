package main

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/auth"
	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/secret"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/server"
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	di.ProvideDeps(ctx)
	config := do.MustInvoke[*utils.Config](nil)

	switch config.SrvType {
	case commonServer.GRPCType:
		client, err := db.NewMongodbClient(ctx, config)
		if err != nil {
			panic(err)
		}
		defer client.Disconnect(ctx)

		encr := do.MustInvoke[encryptor.Encryptor](nil)
		hasher := do.MustInvoke[password.PassowrdHasher](nil)
		repo := do.MustInvoke[domain.Repository](nil)
		userRepo := do.MustInvoke[userDomain.Repository](nil)

		a := app.NewApplication(encr, hasher, repo, userRepo)
		s := server.NewServer(a)
		opts := []grpc.ServerOption{
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.AuthFunc)),
		}
		commonServer.RunGRPCServer(config.GrpcEndpoint, opts, func(srv *grpc.Server) {
			secret.RegisterSecretKeeperServer(srv, s)
			reflection.Register(srv)
		})
	case commonServer.GatewayType:
		commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
			secret.RegisterSecretKeeperHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, opts)
		})
	default:
		panic(errors.UnknownServerType)
	}
}
