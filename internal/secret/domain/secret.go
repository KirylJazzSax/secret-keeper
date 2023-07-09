package domain

import (
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Secret struct {
	Id    int64
	Title string
	Body  string
	User  userDomain.User
}

func NewSecret(title string, body string, user userDomain.User) *Secret {
	return &Secret{
		Title: title,
		Body:  body,
		User:  user,
	}
}
