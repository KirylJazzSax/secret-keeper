package gapi

import (
	"context"
	"secret_keeper/errors"
	"secret_keeper/pb"
)

func (s *Server) DeleteAllSecrets(ctx context.Context, req *pb.DeleteAllSecretsRequest) (*pb.DeleteAllSecretsResponse, error) {
	authPayload, err := s.getAuthPayload(ctx)

	if err != nil {
		errors.LogErr(err)
		return nil, errors.UnAuthErr()
	}

	if err = s.repository.DeleteAllSecrets(authPayload.Email); err != nil {
		return nil, errors.LogErrAndCreateInternal(err)
	}

	return &pb.DeleteAllSecretsResponse{}, nil
}
