package repository

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/secret/domain"
	"go.mongodb.org/mongo-driver/bson"
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
	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	if _, err := coll.InsertOne(ctx, s); err != nil {
		return err
	}
	return nil
}
func (r *MongoRepository) SecretsList(ctx context.Context, userId string) ([]*domain.Secret, error) {
	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	scrs := []*domain.Secret{}
	cursor, err := coll.Find(ctx, bson.D{{"user", id}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &scrs); err != nil {
		return nil, err
	}

	return scrs, nil
}
func (r *MongoRepository) GetSecret(ctx context.Context, id string, userId string) (*domain.Secret, error) {
	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	uId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	sId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := coll.FindOne(ctx, bson.D{{"user", uId}, {"_id", sId}}).Decode(scr); err != nil {
		return nil, err
	}

	return scr, nil
}
func (r *MongoRepository) DeleteSecret(ctx context.Context, id string, userId string) error {
	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	uId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	sId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if _, err := coll.DeleteOne(ctx, bson.D{{"user", uId}, {"_id", sId}}); err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) DeleteAllSecrets(ctx context.Context, userId string) error {
	coll := r.client.Database(db.DB).Collection(db.SecretsCollection)
	uId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	if _, err := coll.DeleteMany(ctx, bson.D{{"user", uId}}); err != nil {
		return err
	}

	return nil
}
