package command

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type ShowPayload struct {
	Id       string
	Email    string
	Password string
	Decoded  *domain.Secret
}

type ShowHandlerType common.CommandHandler[*ShowPayload]

type ShowHandler struct {
	hasher   password.PassowrdHasher
	encr     encryptor.Encryptor
	repo     domain.Repository
	userRepo userDomain.Repository
}

func (h *ShowHandler) Handle(ctx context.Context, p *ShowPayload) error {
	u, err := h.userRepo.GetUser(ctx, p.Email)
	if err != nil {
		return err
	}

	fmt.Println(u.Id.Hex(), p.Id)
	s, err := h.repo.GetSecret(ctx, p.Id, u.Id.Hex())
	if err != nil {
		return err
	}

	var decrypted string
	if err := h.encr.Decrypt(s.Body, &decrypted); err != nil {
		return err
	}

	s.Body = decrypted
	p.Decoded = s

	return nil
}

func NewShowHandler(encr encryptor.Encryptor, hasher password.PassowrdHasher, repo domain.Repository, userRepo userDomain.Repository) *ShowHandler {
	return &ShowHandler{
		hasher:   hasher,
		encr:     encr,
		repo:     repo,
		userRepo: userRepo,
	}
}
