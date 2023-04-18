package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"secret_keeper/errors"
	"secret_keeper/gapi"
	"secret_keeper/internal/di"
	"secret_keeper/pb"
	"secret_keeper/repository"
	"secret_keeper/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func runGatewayServer(fromEndpoint string) {
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

	pb.RegisterSecretKeeperHandlerFromEndpoint(context.Background(), grpcMux, fromEndpoint, opts)
	if err := http.Serve(listener, mux); err != nil {
		panic(err)
	}
}

func main() {
	if err := di.ProvideDeps("."); err != nil {
		panic(err)
	}
	defer do.Shutdown[repository.Repository](nil)

	server := do.MustInvoke[*gapi.Server](nil)

	config := do.MustInvoke[*utils.Config](nil)

	grpcServer := grpc.NewServer()

	endpoint := fmt.Sprintf(":%s", config.Port)
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		errors.LogErr(err)
	}

	pb.RegisterSecretKeeperServer(grpcServer, server)
	reflection.Register(grpcServer)

	go runGatewayServer(endpoint)

	log.Info().Msg("server runs")
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}
