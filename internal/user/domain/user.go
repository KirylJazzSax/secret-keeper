package domain

import (
	"time"
)

type User struct {
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	Password  string    `bson:"password"`
}

func NewUser(email string, createdAt time.Time, password string) *User {
	return &User{
		Email:     email,
		CreatedAt: createdAt,
		Password:  password,
	}
}
