package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	Password  string             `bson:"password"`
}

func NewUser(email string, createdAt time.Time, password string) *User {
	return &User{
		Id:        primitive.NewObjectId(),
		Email:     email,
		CreatedAt: createdAt,
		Password:  password,
	}
}
