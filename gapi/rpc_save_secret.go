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
		return nil, UnAuthErr()
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

	err = s.store.CreateSecret(secret, authPayload.Email)

	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	return &pb.SaveSecretResponse{
		Secret: secret,
	}, nil
}
