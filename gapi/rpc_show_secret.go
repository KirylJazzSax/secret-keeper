package gapi

import (
	"context"
	"secret_keeper/db"
	"secret_keeper/encryptor"
	"secret_keeper/errors"
	"secret_keeper/password"
	"secret_keeper/pb"
)

func (s *Server) ShowSecret(ctx context.Context, req *pb.ShowSecretRequest) (*pb.ShowSecretResponse, error) {
	user, err := s.store.GetUser(req.Email)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.UnAuthErr()
	}

	if err = password.Check(req.Password, user.Password); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	secret, err := s.store.GetSecret(uint64(req.Id), user.Email)

	if err == db.ErrNotExists {
		errors.LogErr(err)
		return nil, errors.ErrNotFound()
	}

	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	if secret.Body, err = encryptor.Decrypt(secret.Body, s.config.SECRET_KEY, s.config.IV); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &pb.ShowSecretResponse{
		Secret: secret,
	}, nil
}
