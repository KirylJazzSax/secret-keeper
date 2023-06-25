package validation

import (
	"errors"
	"net/mail"
)

var InvalidEmail = errors.New("invalid email.")

type Validator interface {
	ValidateEmail(email string) error
}

type SimpleValidator struct {
}

func (v *SimpleValidator) ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return InvalidEmail
	}
	return nil
}

func NewSimpleValidator() *SimpleValidator {
	return &SimpleValidator{}
}
