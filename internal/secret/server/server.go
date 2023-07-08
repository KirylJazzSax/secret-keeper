package server

import (
	"context"
	"go/token"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/secret"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/command"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
)

type Server struct {
	secret.UnimplementedSecretKeeperServer
	application *app.Application
}

func NewServer(application *app.Application) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) SaveSecret(ctx context.Context, r *secret.SaveSecretRequest) (*secret.SaveSecretResponse, error) {
	u := ctx.Value("user").(*token.Payload)

	p := &command.SavePayload{
		Title:  r.Title,
		Body:   r.Body,
		Email:  u.Email,
		Secret: &domain.Secret{},
	}
	if err := s.application.Commands.Save.Handle(ctx, p); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &secret.SaveSecretResponse{
		Secret: p.Secret,
	}, nil
}

func (s *Server) SecretsList(ctx context.Context, r *secret.SecretsListRequest) (*secret.SecretsListResponse, error) {
	return &secret.SecretsListResponse{}, nil
}

func (s *Server) ShowSecret(ctx context.Context, r *secret.ShowSecretRequest) (*secret.ShowSecretResponse, error) {
	return &secret.ShowSecretResponse{
		Secret: &secret.Secret{
			Id:    0,
			Title: "",
			Body:  "",
		},
	}, nil
}

func (s *Server) DeleteSecret(ctx context.Context, r *secret.DeleteSecretRequest) (*secret.DeleteSecretResponse, error) {
	return &secret.DeleteSecretResponse{}, nil
}

func (s *Server) DeleteAllSecrets(ctx context.Context, r *secret.DeleteAllSecretsRequest) (*secret.DeleteAllSecretsResponse, error) {
	return &secret.DeleteAllSecretsResponse{}, nil
}
