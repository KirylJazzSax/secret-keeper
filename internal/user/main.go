package main

import (
	"context"
	"os"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/user"
	"github.com/KirylJazzSax/secret-keeper/internal/common/logs"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	"github.com/KirylJazzSax/secret-keeper/internal/user/server"
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

		hasher := do.MustInvoke[password.PassowrdHasher](nil)
		validator := do.MustInvoke[validation.Validator](nil)
		repo := do.MustInvokeNamed[domain.Repository](nil, "users-repo")

		a := app.NewApplication(validator, hasher, repo)
		s := server.NewServer(a)

		opts := []grpc.ServerOption{
			grpc.UnaryInterceptor(logging.UnaryServerInterceptor(logs.InterceptorLogger(zerolog.New(os.Stdout)))),
		}

		commonServer.RunGRPCServer(config.GrpcEndpoint, opts, func(srv *grpc.Server) {
			user.RegisterUsersServiceServer(srv, s)
			reflection.Register(srv)
		})
	case commonServer.GatewayType:
		commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
			user.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, opts)
		})
	default:
		panic(errors.UnknownServerType)
	}
}
