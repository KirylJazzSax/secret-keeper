package app

import (
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app/command"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app/query"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateUser command.CreateUserHandlerType
}

type Queries struct {
	QueryUser query.QueryUserHandlerType
}

func NewApplication(
	validator validation.Validator,
	hasher password.PassowrdHasher,
	repository domain.Repository,
) *Application {
	return &Application{
		Commands: Commands{
			command.NewCreateUserHandler(validator, hasher, repository),
		},
		Queries: Queries{
			query.NewQueryUserHandler(repository),
		},
	}
}
