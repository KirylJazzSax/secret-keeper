package gapi

import (
	"context"
	"secret_keeper/encryptor"
	"secret_keeper/errors"
	"secret_keeper/pb"

	"github.com/samber/do"
)

func (s *Server) SaveSecret(ctx context.Context, req *pb.SaveSecretRequest) (*pb.SaveSecretResponse, error) {
	authPayload, err := s.getAuthPayload(ctx)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.UnAuthErr()
	}
	encr, err := do.Invoke[encryptor.Encryptor](nil)

	if err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	var encoded string
	if err = encr.Encrypt(req.Body, &encoded); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	secret := &pb.Secret{
		Title: req.Title,
		Body:  encoded,
	}

	if err = s.store.CreateSecret(secret, authPayload.Email); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &pb.SaveSecretResponse{
		Secret: secret,
	}, nil
}
