package domain

import "time"

type User struct {
	Email     string
	CreatedAt time.Time
	Password  string
}

func NewUser(email string, createdAt time.Time, password string) *User {
	return &User{
		Email:     email,
		CreatedAt: createdAt,
		Password:  password,
	}
}
