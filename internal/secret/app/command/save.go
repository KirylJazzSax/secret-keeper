package command

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common"
)

type Payload struct {
}

type SaveHandlerType common.CommandHandler[*Payload]

type SaveHandler struct {
}

func (h *SaveHandlerType) Handle(ctx context.Context, p *Payload) error {
	return nil
}

func NewSaveHandler() *SaveHandler {
	return &SaveHandler{}
}
