package repository

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
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
	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	if _, err := coll.InsertOne(ctx, s); err != nil {
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
