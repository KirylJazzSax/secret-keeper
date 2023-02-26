package gapi

import (
	"fmt"
	"secret_keeper/db"
	"secret_keeper/pb"
	"secret_keeper/token"
	"secret_keeper/utils"

	"github.com/samber/do"
)

type Server struct {
	pb.UnimplementedSecretKeeperServer
	store        db.Store
	tokenManager token.Maker
	config       *utils.Config
}

func NewServer(store db.Store, config *utils.Config) (*Server, error) {
	manager, err := do.Invoke[token.Maker](nil)
	if err != nil {
		return nil, fmt.Errorf("could not create tokenmanager")
	}

	return &Server{
		store:        store,
		tokenManager: manager,
		config:       config,
	}, nil
}
