package common

import "context"

type QueryHandler[T any, R any] interface {
	Query(ctx context.Context, payload T) (R, error)
}
