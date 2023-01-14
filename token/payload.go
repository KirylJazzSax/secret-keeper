package token

import (
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	Email     string                `json:"email"`
	IssuedAt  timestamppb.Timestamp `json:"issued_at"`
	ExpiredAt timestamppb.Timestamp `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(email string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		Email:     email,
		IssuedAt:  *timestamppb.Now(),
		ExpiredAt: *timestamppb.New(time.Now().Add(duration)),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt.AsTime()) {
		return ErrExpiredToken
	}
	return nil
}
