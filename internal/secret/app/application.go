package app

import (
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/command"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/query"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Commands struct {
	Save command.SaveHandlerType
}

type Queries struct {
	All query.AllQueryHanlderType
}

type Application struct {
	Commands Commands
	Queries  Queries
}

func NewApplication(encr encryptor.Encryptor, repo domain.Repository, userRepo userDomain.Repository) *Application {
	return &Application{
		Commands: Commands{
			Save: command.NewSaveHandler(encr, repo, userRepo),
		},
		Queries: Queries{
			All: query.NewAllQueryHandler(repo),
		},
	}
}
