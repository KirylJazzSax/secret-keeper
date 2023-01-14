package gapi

import (
	"context"
	"secret_keeper/errors"
	"secret_keeper/pb"
)

func (s *Server) SecretsList(ctx context.Context, req *pb.SecretsListRequest) (*pb.SecretsListResponse, error) {
	authPayload, err := s.getAuthPayload(ctx)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	list, err := s.store.SecretsList(authPayload.Email)
	if err != nil {
		errors.LogErr(err)
		return nil, errors.ErrInternal()
	}

	return &pb.SecretsListResponse{
		Secrets: list,
	}, nil
}
