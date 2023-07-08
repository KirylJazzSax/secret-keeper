package domain

import (
	userDomain "github.com/KirylJazzSax/secret-keeper/internal/user/domain"
)

type Secret struct {
	Title string          `bson:"title"`
	Body  string          `bson:"body"`
	User  userDomain.User `bson:"inline"`
}

func NewSecret(title string, body string, user userDomain.User) *Secret {
	return &Secret{
		Title: title,
		Body:  body,
		User:  user,
	}
}
