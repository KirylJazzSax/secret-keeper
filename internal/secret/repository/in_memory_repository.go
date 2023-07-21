package repository

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
)

type InMemoryRepository struct {
	secrets map[string]*domain.Secret
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		secrets: make(map[string]*domain.Secret),
	}
}

func (r *InMemoryRepository) CreateSecret(ctx context.Context, s *domain.Secret) error {
	r.secrets[s.Id.Hex()] = s
	return nil
}
func (r *InMemoryRepository) SecretsList(ctx context.Context, userId string) ([]*domain.Secret, error) {
	scrs := []*domain.Secret{}

	for _, v := range r.secrets {
		if v.User.Hex() == userId {
			scrs = append(scrs, v)
		}
	}

	return scrs, nil
}
func (r *InMemoryRepository) GetSecret(ctx context.Context, id string, userId string) (*domain.Secret, error) {
	s, ok := r.secrets[id]
	if !ok {
		return nil, errors.ErrNotExists
	}

	if s.User.Hex() != userId {
		return nil, errors.ErrNotExists
	}

	return s, nil
}

func (r *InMemoryRepository) DeleteSecret(ctx context.Context, id string, userId string) error {
	s, err := r.GetSecret(ctx, id, userId)
	if err != nil {
		return err
	}

	delete(r.secrets, s.Id.Hex())

	return nil
}

func (r *InMemoryRepository) DeleteAllSecrets(ctx context.Context, userId string) error {
	for _, v := range r.secrets {
		if v.User.Hex() == userId {
			delete(r.secrets, v.Id.Hex())
		}
	}

	return nil
}
