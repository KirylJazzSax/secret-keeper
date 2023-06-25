package command

import (
	"context"
	"time"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Payload struct {
	Email    string
	Password string
}

type CreateUserHandlerType common.CommandHandler[*Payload]

type CreateUserHandler struct {
	validator  validation.Validator
	hasher     password.PassowrdHasher
	repository domain.Repository
}

func NewCreateUserHandler(
	validator validation.Validator,
	hasher password.PassowrdHasher,
	repository domain.Repository,
) *CreateUserHandler {
	return &CreateUserHandler{
		validator:  validator,
		hasher:     hasher,
		repository: repository,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, payload *Payload) error {
	if err := h.validator.ValidateEmail(payload.Email); err != nil {
		return err
	}

	var hash string
	if err := h.hasher.Hash(payload.Password, &hash); err != nil {
		return err
	}

	if err := h.repository.CreateUser(context.Background(), domain.NewUser(payload.Email, time.Now(), hash)); err != nil {
		return err
	}
	return nil
}
