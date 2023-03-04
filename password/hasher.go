package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PassowrdHasher interface {
	Hash(password string, hashTo *string) error
	Check(password, hash string) error
}

type SimplePasswordHasher struct {
}

func (h *SimplePasswordHasher) Hash(password string, hashTo *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("hash password: %s", err)
	}

	*hashTo = string(hash)
	return nil
}

func (h *SimplePasswordHasher) Check(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewSimplePasswordHasher() *SimplePasswordHasher {
	return &SimplePasswordHasher{}
}
