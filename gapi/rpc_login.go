package gapi

import (
	"context"
	"secret_keeper/errors"
	"secret_keeper/password"
	"secret_keeper/pb"
	"secret_keeper/repository"

	"github.com/samber/do"
)

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.repository.GetUser(request.Email)

	if err == repository.ErrNotExists {
		errors.LogErr(err)
		return nil, errors.ErrNotFound()
	}

	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	hasher, err := do.Invoke[password.PassowrdHasher](nil)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	if err = hasher.Check(request.Password, user.Password); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	token, payload, err := server.tokenManager.CreateToken(request.Email, server.config.AccessTokenDuration)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &pb.LoginResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: &payload.ExpiredAt,
	}, nil
}
