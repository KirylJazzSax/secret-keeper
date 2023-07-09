package command

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type SavePayload struct {
	Title  string
	Body   string
	Email  string
	Secret *domain.Secret
}

type SaveHandlerType common.CommandHandler[*SavePayload]

type SaveHandler struct {
	encr     encryptor.Encryptor
	repo     domain.Repository
	userRepo userDomain.Repository
}

func (h *SaveHandler) Handle(ctx context.Context, p *SavePayload) error {
	u, err := h.userRepo.GetUser(ctx, p.Email)
	if err != nil {
		return err
	}

	var encoded string
	if err := h.encr.Encrypt(p.Body, &encoded); err != nil {
		return err
	}

	s := domain.NewSecret(p.Title, encoded, *u)

	if err := h.repo.CreateSecret(ctx, s); err != nil {
		return err
	}

	p.Secret = s

	return nil
}

func NewSaveHandler(encr encryptor.Encryptor, repo domain.Repository, userRepo userDomain.Repository) *SaveHandler {
	return &SaveHandler{
		encr:     encr,
		repo:     repo,
		userRepo: userRepo,
	}
}
