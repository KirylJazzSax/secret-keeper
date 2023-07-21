package repository

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type InMemoryUserRepository struct {
	users map[string]*domain.User
}

func (r *InMemoryUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	r.users[u.Email] = u
	return nil
}

func (r *InMemoryUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	return r.users[email], nil
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}
