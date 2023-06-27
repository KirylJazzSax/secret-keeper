package server

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/auth/app"
	"github.com/KirylJazzSax/secret-keeper/internal/auth/command"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/auth"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	application app.Application
}

func NewServer(application app.Application) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) LoginUser(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	payload := &command.Payload{
		Email:     request.Email,
		Password:  request.Password,
		Token:     "",
		ExpiredAt: timestamppb.Now(),
	}
	user, err := s.application.LoginUserCommand.Handle(ctx, payload)

	if err == errors.ErrNotExists {
		errors.LogErr(err)
		return nil, errors.ErrNotFound()
	}

	if err == errors.ErrEmailOrPasswordNotValid {
		var violations []*errdetails.BadRequest_FieldViolation
		violations = append(violations, &errdetails.BadRequest_FieldViolation{Field: "email", Description: err.Error()})
		return nil, errors.InvalidArgumentError(violations)
	}

	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &auth.LoginResponse{
		AccessToken:          payload.Token,
		AccessTokenExpiresAt: payload.ExpiredAt,
	}, nil
}
