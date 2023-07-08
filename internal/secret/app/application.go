package app

import (
	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/command"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
)

type Commands struct {
	Save command.SaveHandlerType
}

type Application struct {
	Commands Commands
}

func NewApplication(encr encryptor.Encryptor, repo domain.Repository) *Application {
	return &Application{
		Commands: Commands{
			Save: command.NewSaveHandler(encr, repo),
		},
	}
}
