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

func RunGRPCServer(endpoint string, cb func(s *grpc.Server)) error {
	grpcServer := grpc.NewServer()

	cb(grpcServer)

	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}

	log.Info().Msg("server runs")
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func RunGatewayServer(corsOrigin string, httpPort string, cb func(mux *runtime.ServeMux, opts []grpc.DialOption)) error {
	grpcMux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")

		grpcMux.ServeHTTP(w, r)
	}))

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", httpPort))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cb(mux, opts)
	if err := http.Serve(listener, mux); err != nil {
		return err
	}

	return nil
}
