package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const (
	GRPCType    = "grpc"
	GatewayType = "gateway"
)

func RunGRPCServer(endpoint string, opts []grpc.ServerOption, cb func(s *grpc.Server)) error {
	grpcServer := grpc.NewServer(opts...)

	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}

	cb(grpcServer)

	log.Info().Msg("server runs")
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func RunGatewayServer(corsOrigin string, httpPort string, cb func(mux *runtime.ServeMux, opts []grpc.DialOption)) error {
	mux := runtime.NewServeMux()
	httpMux := http.NewServeMux()

	httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", httpPort))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cb(mux, opts)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{corsOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "ResponseType"},
	})

	handler := c.Handler(httpMux)
	log.Info().Msg("running http.")
	if err := http.Serve(listener, handler); err != nil {
		return err
	}

	return nil
}
