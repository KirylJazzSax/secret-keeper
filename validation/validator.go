package validation

import (
	"fmt"
	"net/mail"
)

type Validator interface {
	ValidateEmail(email string) error
}

type SimpleValidator struct {
}

func (v *SimpleValidator) ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email")
	}
	return nil
}

func NewSimpleValidator() *SimpleValidator {
	return &SimpleValidator{}
}
