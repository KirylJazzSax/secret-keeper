package repository

import "secret_keeper/pb"

type Repository interface {
	CreateUser(user *User) error
	GetUser(email string) (*User, error)
	CreateSecret(secret *pb.Secret, email string) error
	SecretsList(email string) ([]*pb.Secret, error)
	GetSecret(id uint64, email string) (*pb.Secret, error)
	DeleteSecret(id uint64, email string) error
	DeleteAllSecrets(email string) error
}
