package domain

import (
	"context"
)

type Repository interface {
	CreateSecret(ctx context.Context, s *Secret) error
	SecretsList(ctx context.Context, email string) ([]*Secret, error)
	GetSecret(ctx context.Context, id uint64, email string) (*Secret, error)
	DeleteSecret(ctx context.Context, id uint64, email string) error
	DeleteAllSecrets(ctx context.Context, email string) error
}
