package domain

import (
	"context"
)

type Repository interface {
	CreateSecret(ctx context.Context, s *Secret) error
	SecretsList(ctx context.Context, userId string) ([]*Secret, error)
	GetSecret(ctx context.Context, id string, userId string) (*Secret, error)
	DeleteSecret(ctx context.Context, id string, userId string) error
	DeleteAllSecrets(ctx context.Context, userId string) error
}
