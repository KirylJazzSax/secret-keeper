package command

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Payload struct {
	Email     string
	Password  string
	AuthToken string
	ExpiredAt timestamppb.Timestamp
}

type LoginUserHandlerType common.CommandHandler[*Payload]

type LoginUserHandler struct {
	tokenManager token.Maker
	hasher       password.PassowrdHasher
	repo         domain.Repository
	config       *utils.Config
}

func (h *LoginUserHandler) Handle(ctx context.Context, payload *Payload) error {
	user, err := h.repo.GetUser(ctx, payload.Email)

	if err == errors.ErrNotExists {
		return errors.ErrNotExists
	}

	if err != nil {
		return errors.ErrInteralServer
	}

	if err = h.hasher.Check(payload.Password, user.Password); err != nil {
		return errors.ErrEmailOrPasswordNotValid
	}

	token, p, err := h.tokenManager.CreateToken(payload.Email, h.config.AccessTokenDuration)
	if err != nil {
		return errors.ErrInteralServer
	}

	payload.AuthToken = token
	payload.ExpiredAt = p.ExpiredAt

	return nil
}

func NewLoginHandler(
	tokenManager token.Maker,
	hasher password.PassowrdHasher,
	repo domain.Repository,
	config *utils.Config,
) *LoginUserHandler {
	return &LoginUserHandler{
		tokenManager: tokenManager,
		hasher:       hasher,
		repo:         repo,
		config:       config,
	}
}
