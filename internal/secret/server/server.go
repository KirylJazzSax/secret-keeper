package server

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/secret"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app"
)

type Server struct {
	secret.UnimplementedUsersServiceServer
	application *app.Application
}

func NewServer(application *app.Application) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) SaveSecret(ctx context.Context, r *secret.SaveSecretRequest) (*secret.SaveSecretResponse, error) {
	return &secret.SaveSecretResponse{
		Secret: &secret.Secret{
			Id:    0,
			Title: "",
			Body:  "",
		},
	}, nil
}
