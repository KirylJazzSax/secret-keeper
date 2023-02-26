package gapi

import (
	"context"
	"secret_keeper/encryptor"
	"secret_keeper/errors"
	"secret_keeper/pb"
)

func (s *Server) SaveSecret(ctx context.Context, req *pb.SaveSecretRequest) (*pb.SaveSecretResponse, error) {
	authPayload, err := s.getAuthPayload(ctx)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.UnAuthErr()
	}

	str, err := encryptor.Encrypt(req.Body, s.config.SECRET_KEY, s.config.IV)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	secret := &pb.Secret{
		Title: req.Title,
		Body:  str,
	}

	if err = s.store.CreateSecret(secret, authPayload.Email); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &pb.SaveSecretResponse{
		Secret: secret,
	}, nil
}
