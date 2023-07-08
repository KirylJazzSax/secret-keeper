package app

import (
	"github.com/KirylJazzSax/secret-keeper/internal/secret/app/command"
)

type Commands struct {
	Save command.SaveHandlerType
}

type Application struct {
	Commands Commands
}

func NewApplication() *Application {
	return &Application{
		Commands: Commands{
			Save: command.NewSaveHandler(),
		},
	}
}
