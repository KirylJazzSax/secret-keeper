package gapi

import (
	"secret_keeper/pb"
	"secret_keeper/repository"
	"secret_keeper/token"
	"secret_keeper/utils"
)

type Server struct {
	pb.UnimplementedSecretKeeperServer
	repository   repository.Repository
	tokenManager token.Maker
	config       *utils.Config
}

func NewServer(repo repository.Repository, manager token.Maker, config *utils.Config) *Server {
	return &Server{
		repository:   repo,
		tokenManager: manager,
		config:       config,
	}
}
