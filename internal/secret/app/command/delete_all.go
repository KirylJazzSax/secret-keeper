package command

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
)

type DeleteAllPayload struct {
	UserId string
}

type DeleteAllHandlerType common.CommandHandler[*DeleteAllPayload]

type DeleteAllHandler struct {
	repo domain.Repository
}

func (h *DeleteAllHandler) Handle(ctx context.Context, p *DeleteAllPayload) error {
	return h.repo.DeleteAllSecrets(ctx, p.UserId)
}

func NewDeleteHandler(repo domain.Repository) *DeleteAllHandler {
	return &DeleteAllHandler{
		repo: repo,
	}
}
