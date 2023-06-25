package server

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/gen/user"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app"
	"github.com/KirylJazzSax/secret-keeper/internal/user/app/command"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	user.UnimplementedUsersServiceServer
	application *app.Application
}

func NewServer(application *app.Application) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) CreateUser(ctx context.Context, request *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	err := s.application.Commands.CreateUser.Handle(ctx, &command.Payload{
		Email:    request.Email,
		Password: request.Password,
	})

	if err == validation.InvalidEmail {
		var violations []*errdetails.BadRequest_FieldViolation
		violations = append(violations, &errdetails.BadRequest_FieldViolation{Field: "email", Description: err.Error()})
		return nil, errors.InvalidArgumentError(violations)
	}

	if err == errors.ErrExists {
		var violations []*errdetails.BadRequest_FieldViolation
		violations = append(violations, &errdetails.BadRequest_FieldViolation{Field: "email", Description: err.Error()})
		return nil, errors.InvalidArgumentError(violations)
	}

	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &user.CreateUserResponse{
		User: &user.User{
			Email:     request.Email,
			CreatedAt: timestamppb.Now(),
		},
	}, nil
}
