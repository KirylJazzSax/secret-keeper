package common

import "context"

type CommandHandler[T any] interface {
	Handle(ctx context.Context, c T) error
}
