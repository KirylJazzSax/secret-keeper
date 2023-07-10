package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Secret struct {
	Id    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title"`
	Body  string             `bson:"body"`
	User  primitive.ObjectID `bson:"user"`
}

func NewSecret(title string, body string, user primitive.ObjectID) *Secret {
	return &Secret{
		Title: title,
		Body:  body,
		User:  user,
	}
}
