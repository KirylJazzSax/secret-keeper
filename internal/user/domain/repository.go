package domain

import "context"

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUser(ctx context.Context, email string) (*User, error)
}
