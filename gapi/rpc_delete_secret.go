package gapi

import (
	"context"

	"secret_keeper/errors"
	"secret_keeper/pb"
)

func (s *Server) DeleteSecret(ctx context.Context, req *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	authPayload, err := s.getAuthPayload(ctx)

	if err != nil {
		errors.LogErr(err)
		return nil, errors.UnAuthErr()
	}

	if err = s.repository.DeleteSecret(uint64(req.Id), authPayload.Email); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &pb.DeleteSecretResponse{}, nil
}
