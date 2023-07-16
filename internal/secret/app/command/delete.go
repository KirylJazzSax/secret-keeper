package command

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
)

type DeletePayload struct {
	Id     string
	UserId string
}

type DeleteHandlerType common.CommandHandler[*DeletePayload]

type DeleteHandler struct {
	repo domain.Repository
}

func (h *DeleteHandler) Handle(ctx context.Context, p *DeletePayload) error {
	return h.repo.DeleteSecret(ctx, p.Id, p.UserId)
}

func NewDeleteHandler(repo domain.Repository) *DeleteHandler {
	return &DeleteHandler{
		repo: repo,
	}
}
