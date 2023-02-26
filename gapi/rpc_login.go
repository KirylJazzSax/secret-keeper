package gapi

import (
	"context"
	"secret_keeper/db"
	"secret_keeper/errors"
	"secret_keeper/password"
	"secret_keeper/pb"
)

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.store.GetUser(request.Email)

	if err == db.ErrNotExists {
		errors.LogErr(err)
		return nil, errors.ErrNotFound()
	}

	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	if err = password.Check(request.Password, user.Password); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	token, payload, err := server.tokenManager.CreateToken(request.Email, server.config.AccessTokenDuration)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	return &pb.LoginResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: &payload.ExpiredAt,
	}, nil
}
