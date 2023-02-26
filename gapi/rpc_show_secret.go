package gapi

import (
	"context"
	"secret_keeper/db"
	"secret_keeper/encryptor"
	"secret_keeper/errors"
	"secret_keeper/password"
	"secret_keeper/pb"

	"github.com/samber/do"
)

func (s *Server) ShowSecret(ctx context.Context, req *pb.ShowSecretRequest) (*pb.ShowSecretResponse, error) {
	user, err := s.store.GetUser(req.Email)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.UnAuthErr()
	}

	hasher, err := do.Invoke[password.PassowrdHasher](nil)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	if err = hasher.Check(req.Password, user.Password); err != nil {
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

	encr, err := do.Invoke[encryptor.Encryptor](nil)
	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	var decripted string
	if err = encr.Decrypt(secret.Body, &decripted); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	secret.Body = decripted

	return &pb.ShowSecretResponse{
		Secret: secret,
	}, nil
}
