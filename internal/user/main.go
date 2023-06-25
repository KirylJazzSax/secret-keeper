package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/KirylJazzSax/secret-keeper/internal/common/di"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/user"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app"
	"github.com/KirylJazzSax/secret-keeper/internal/user/repository"
	"github.com/KirylJazzSax/secret-keeper/internal/user/server"

	"github.com/boltdb/bolt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	grpcServer := grpc.NewServer()

	endpoint := fmt.Sprintf(":%s", config.Port)
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		panic(err)
	}

	user.RegisterUsersServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	go runGatewayServer(endpoint)

	log.Info().Msg("server runs")
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}

}

func runGatewayServer(fromEndpoint string) error {
	config := do.MustInvoke[*utils.Config](nil)

	grpcMux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.Cors)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")

		grpcMux.ServeHTTP(w, r)
	}))

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", config.HttpPort))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	docHandler := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", docHandler))

	user.RegisterUsersServiceHandlerFromEndpoint(context.Background(), grpcMux, fromEndpoint, opts)
	if err := http.Serve(listener, mux); err != nil {
		panic(err)
	}

	return nil
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
