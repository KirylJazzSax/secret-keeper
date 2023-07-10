package query

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/secret/common"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
)

type AllPayload struct {
	UserId string
}

type AllQueryHanlderType common.QueryHandler[AllPayload, []*domain.Secret]

type AllQueryHanlder struct {
	repo domain.Repository
}

func NewAllQueryHandler(repo domain.Repository) *AllQueryHanlder {
	return &AllQueryHanlder{
		repo: repo,
	}
}

func (h *AllQueryHanlder) Handle(ctx context.Context, p *AllPayload) ([]*domain.Secret, error) {
	return h.repo.SecretsList(ctx, p.UserId)
}
