package main

import (
	"context"
	"os"

	"github.com/KirylJazzSax/secret-keeper/internal/common/auth"
	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/secret"
	"github.com/KirylJazzSax/secret-keeper/internal/common/logs"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/server"
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
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
		do.MustInvoke[*db.Db](nil)
		defer do.Shutdown[*db.Db](nil)

		encr := do.MustInvoke[encryptor.Encryptor](nil)
		hasher := do.MustInvoke[password.PassowrdHasher](nil)
		repo := do.MustInvoke[domain.Repository](nil)
		userRepo := do.MustInvoke[userDomain.Repository](nil)

		a := app.NewApplication(encr, hasher, repo, userRepo)
		s := server.NewServer(a)
		opts := []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				grpc_auth.UnaryServerInterceptor(auth.AuthFunc),
				logging.UnaryServerInterceptor(logs.InterceptorLogger(zerolog.New(os.Stdout))),
			),
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
