package gapi

import (
	"context"

	"secret_keeper/errors"
	"secret_keeper/password"
	"secret_keeper/pb"
	"secret_keeper/repository"
	"secret_keeper/validation"

	"github.com/samber/do"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &pb.User{
		Email:     request.Email,
		CreatedAt: timestamppb.Now(),
	}

	v, err := do.Invoke[validation.Validator](nil)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	if violations := validateRequest(user, v); len(violations) > 0 {
		return nil, errors.InvalidArgumentError(violations)
	}

	hasher, err := do.Invoke[password.PassowrdHasher](nil)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	var hash string
	if err = hasher.Hash(request.Password, &hash); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	dbUser := &repository.User{User: user, Password: hash}
	if err = server.repository.CreateUser(dbUser); err != nil {
		if err == repository.ErrExists {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	return &pb.CreateUserResponse{User: user}, nil
}

func validateRequest(user *pb.User, v validation.Validator) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := v.ValidateEmail(user.Email); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{Field: "email", Description: err.Error()})
	}
	return violations
}
