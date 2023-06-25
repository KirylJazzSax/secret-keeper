package app

import (
	"github.com/KirylJazzSax/secret-keeper/internal/auth/command"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Application struct {
	LoginUserCommand command.LoginUserHandlerType
}

func NewApplication(
	tokenManager token.Maker,
	hasher password.PassowrdHasher,
	repo domain.Repository,
	config *utils.Config,
) *Application {
	return &Application{
		LoginUserCommand: command.NewLoginHandler(tokenManager, hasher, repo, config),
	}
}
