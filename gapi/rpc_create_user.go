package gapi

import (
	"context"
	"secret_keeper/db"
	"secret_keeper/password"
	"secret_keeper/pb"
	"secret_keeper/validation"

	"secret_keeper/errors"

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

	violations := validateRequest(user)

	if len(violations) > 0 {
		return nil, errors.InvalidArgumentError(violations)
	}

	hash, err := password.Hash(request.Password)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	dbUser := &db.User{User: user, Password: hash}
	err = server.store.CreateUser(dbUser)
	if err == db.ErrExists {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	return &pb.CreateUserResponse{User: user}, nil
}

func validateRequest(user *pb.User) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateEmail(user.Email); err != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{Field: "email", Description: err.Error()})
	}
	return violations
}
