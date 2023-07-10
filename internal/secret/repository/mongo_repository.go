package repository

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		client: client,
	}
}

func (r *MongoRepository) CreateSecret(ctx context.Context, s *domain.Secret) error {
	userId, err := primitive.ObjectIDFromHex(fmt.Sprintf("%x", s.User.Id))

	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	if _, err := coll.InsertOne(ctx, &SecretDto{
		Title: s.Title,
		Body:  s.Body,
		User:  s.User.Id,
	}); err != nil {
		return err
	}
	return nil
}
func (r *MongoRepository) SecretsList(ctx context.Context, email string) ([]*domain.Secret, error) {
	return []*domain.Secret{}, nil
}
func (r *MongoRepository) GetSecret(ctx context.Context, id uint64, email string) (*domain.Secret, error) {
	return &domain.Secret{}, nil
}
func (r *MongoRepository) DeleteSecret(ctx context.Context, id uint64, email string) error {
	return nil
}
func (r *MongoRepository) DeleteAllSecrets(ctx context.Context, email string) error {
	return nil
}
