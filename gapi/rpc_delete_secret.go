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
		return nil, UnAuthErr()
	}

	s.store.DeleteSecret(uint64(req.Id), authPayload.Email)

	return &pb.DeleteSecretResponse{}, nil
}
