package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const (
	GRPCType    = "grpc"
	GatewayType = "gateway"
)

func RunGRPCServer(endpoint string, cb func(s *grpc.Server)) error {
	grpcServer := grpc.NewServer()

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
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")

		mux.ServeHTTP(w, r)
	}))

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", httpPort))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cb(mux, opts)
	log.Info().Msg("running http.")
	if err := http.Serve(listener, httpMux); err != nil {
		return err
	}

	return nil
}
