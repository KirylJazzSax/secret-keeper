package repository

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type InMemoryUserRepository struct {
	users map[string]*domain.User
}

func (r *InMemoryUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	if _, ok := r.users[u.Email]; ok {
		return errors.ErrExists
	}

	r.users[u.Email] = u
	return nil
}

func (r *InMemoryUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	u, ok := r.users[email]

	if !ok {
		return nil, errors.ErrNotExists
	}

	return u, nil
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}
