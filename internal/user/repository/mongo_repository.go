package repository

import (
	"context"
	"secret-keeper/internal/user/domain"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	client *mongo.Client
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	coll := r.client.Database(db.DB).Collection(db.UsersCollection)
	_, err := coll.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return errors.ErrExists
	}
	return err
}

func (r *MongoUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	coll := r.client.Database(db.DB).Collection(db.UsersCollection)
	filter := bson.D{{"email", email}}
	user := &domain.User{}

	err := coll.FindOne(ctx, filter).Decode(user)

	return user, err
}

func NewMongoUserRepository(client *mongo.Client) *MongoUserRepository {
	return &MongoUserRepository{
		client: client,
	}
}
