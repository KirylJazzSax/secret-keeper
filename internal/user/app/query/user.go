package query

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Payload struct {
	Email string
}

type QueryUserHandlerType common.QueryHandler[*Payload, *domain.User]

type Handler struct {
	repository domain.Repository
}

func NewQueryUserHandler(repository domain.Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (q *Handler) Query(ctx context.Context, p *Payload) (*domain.User, error) {
	return q.repository.GetUser(ctx, p.Email)
}
