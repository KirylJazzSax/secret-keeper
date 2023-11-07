package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	authApp "github.com/KirylJazzSax/secret-keeper/internal/auth/app"
	authServer "github.com/KirylJazzSax/secret-keeper/internal/auth/server"
	commonAuth "github.com/KirylJazzSax/secret-keeper/internal/common/auth"
	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/auth"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/secret"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/user"
	"github.com/KirylJazzSax/secret-keeper/internal/common/logs"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	commonServer "github.com/KirylJazzSax/secret-keeper/internal/common/server"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	secretApp "github.com/KirylJazzSax/secret-keeper/internal/secret/app"
	secretDomain "github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	secretRepo "github.com/KirylJazzSax/secret-keeper/internal/secret/repository"
	secretServer "github.com/KirylJazzSax/secret-keeper/internal/secret/server"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	"github.com/KirylJazzSax/secret-keeper/internal/user/repository"
	"github.com/KirylJazzSax/secret-keeper/internal/user/server"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	os.Setenv("ACCESS_TOKEN_DURATION", "30m")
	os.Setenv("IV", "aaaaaaaaaAAAAAAA")
	os.Setenv("SECRET_KEY", "AaaaaaAAAaaaaaaaaaaAaAaaaaaaaaaa")
	os.Setenv("SYMMETRIC_KEY", "aAaAAAaaaaaaaaaaaaaAaAaaaaaaaaaa")

	ctx := context.Background()
	di.ProvideDeps(ctx)
	provideDeps()

	config := do.MustInvokeNamed[*utils.Config](nil, "users-config")

	tokenManager := do.MustInvoke[token.Maker](nil)
	hasher := do.MustInvoke[password.PassowrdHasher](nil)
	encr := do.MustInvoke[encryptor.Encryptor](nil)
	validator := do.MustInvoke[validation.Validator](nil)
	repo := do.MustInvokeNamed[domain.Repository](nil, "users-repo")
	secretsRepo := do.MustInvokeNamed[secretDomain.Repository](nil, "secrets-repo")

	userApp := app.NewApplication(validator, hasher, repo)
	userServer := server.NewServer(userApp)
	logger := zerolog.New(os.Stdout)

	secretApp := secretApp.NewApplication(encr, hasher, secretsRepo, repo)
	secretServer := secretServer.NewServer(secretApp)
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_auth.UnaryServerInterceptor(commonAuth.AuthFunc),
			logging.UnaryServerInterceptor(logs.InterceptorLogger(logger)),
		),
	}

	authApp := authApp.NewApplication(
		tokenManager,
		hasher,
		repo,
		config,
	)
	authServer := authServer.NewServer(authApp)

	go commonServer.RunGRPCServer(config.GrpcEndpoint, opts, func(srv *grpc.Server) {
		auth.RegisterAuthServiceServer(srv, authServer)
		user.RegisterUsersServiceServer(srv, userServer)
		secret.RegisterSecretKeeperServer(srv, secretServer)
		reflection.Register(srv)
	})

	go commonServer.RunGatewayServer(config.Cors, config.HttpPort, func(mux *runtime.ServeMux, opts []grpc.DialOption) {
		auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, opts)
		user.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, opts)
		secret.RegisterSecretKeeperHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, opts)
	})

	exit := make(chan os.Signal, 1)
	signal.Notify(
		exit,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	<-exit
}

func provideDeps() error {
	dur, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return err
	}

	do.ProvideNamedValue(nil, "users-config", &utils.Config{
		GrpcEndpoint:        ":8000",
		HttpPort:            "8080",
		Cors:                "*",
		IV:                  os.Getenv("IV"),
		SymmetricKey:        os.Getenv("SYMMETRIC_KEY"),
		SecretKey:           os.Getenv("SECRET_KEY"),
		AccessTokenDuration: dur,
	})

	do.ProvideNamedValue(nil, "secrets-config", &utils.Config{
		GrpcEndpoint:        ":8002",
		HttpPort:            "8082",
		Cors:                "*",
		IV:                  os.Getenv("IV"),
		SymmetricKey:        os.Getenv("SYMMETRIC_KEY"),
		SecretKey:           os.Getenv("SECRET_KEY"),
		AccessTokenDuration: dur,
	})

	do.ProvideNamedValue(nil, "auth-config", &utils.Config{
		GrpcEndpoint:        ":8001",
		HttpPort:            "8081",
		Cors:                "*",
		IV:                  os.Getenv("IV"),
		SymmetricKey:        os.Getenv("SYMMETRIC_KEY"),
		SecretKey:           os.Getenv("SECRET_KEY"),
		AccessTokenDuration: dur,
	})

	do.OverrideNamed[domain.Repository](nil, "users-repo", func(i *do.Injector) (domain.Repository, error) {
		return repository.NewInMemoryUserRepository(), nil
	})

	do.OverrideNamed[secretDomain.Repository](nil, "secrets-repo", func(i *do.Injector) (secretDomain.Repository, error) {
		return secretRepo.NewInMemoryRepository(), nil
	})

	return nil
}
