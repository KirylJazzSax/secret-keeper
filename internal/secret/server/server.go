package server

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/secret"
	tokenMaker "github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/command"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/query"
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
	u := ctx.Value("user").(*tokenMaker.Payload)

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
		Secret: &secret.Secret{
			Id:    p.Secret.Id.String(),
			Title: p.Secret.Title,
			Body:  p.Secret.Body,
		},
	}, nil
}

func (s *Server) SecretsList(ctx context.Context, r *secret.SecretsListRequest) (*secret.SecretsListResponse, error) {
	u := ctx.Value("user").(*tokenMaker.Payload)

	p := &query.AllPayload{
		UserId: u.Id,
	}

	secrets, err := s.application.Queries.All.Query(ctx, p)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	resSecrets := make([]*secret.Secret, len(secrets))
	for i, s := range secrets {
		resSecrets[i] = &secret.Secret{
			Id:    s.Id.String(),
			Title: s.Title,
			Body:  s.Body,
		}
	}

	return &secret.SecretsListResponse{
		Secrets: resSecrets,
	}, nil
}

func (s *Server) ShowSecret(ctx context.Context, r *secret.ShowSecretRequest) (*secret.ShowSecretResponse, error) {
	p := &command.ShowPayload{
		Id:       r.Id,
		Email:    r.Email,
		Password: r.Password,
		Decoded:  &domain.Secret{},
	}

	if err := s.application.Commands.Show.Handle(ctx, p); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &secret.ShowSecretResponse{
		Secret: &secret.Secret{
			Id:    p.Decoded.Id.Hex(),
			Title: p.Decoded.Title,
			Body:  p.Decoded.Body,
		},
	}, nil
}

func (s *Server) DeleteSecret(ctx context.Context, r *secret.DeleteSecretRequest) (*secret.DeleteSecretResponse, error) {
	u := ctx.Value("user").(*tokenMaker.Payload)

	p := &command.DeletePayload{
		Id:     r.Id,
		UserId: u.Id,
	}

	if err := *s.application.Commands.Delete.Handle(ctx, p); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &secret.DeleteSecretResponse{}, nil
}

func (s *Server) DeleteAllSecrets(ctx context.Context, r *secret.DeleteAllSecretsRequest) (*secret.DeleteAllSecretsResponse, error) {
	u := ctx.Value("user").(*tokenMaker.Payload)

	p := &command.DeletePayload{
		UserId: u.Id,
	}

	if err := *s.application.Commands.DeleteAll.Handle(ctx, p); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}
	return &secret.DeleteAllSecretsResponse{}, nil
}
