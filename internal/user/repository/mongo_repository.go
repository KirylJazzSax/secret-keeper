package repository

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	client *mongo.Client
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	coll := r.client.Database("secret-keeper").Collection("users")
	_, err := coll.InsertOne(ctx, u)
	return err
}

func (r *MongoUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	coll := r.client.Database("secret-keeper").Collection("users")
	filter := bson.D{{"email", email}}
	var user *domain.User

	err := coll.FindOne(ctx, filter).Decode(user)

	return user, err
}

func NewMongoUserRepository(client *mongo.Client) *MongoUserRepository {
	return &MongoUserRepository{
		client: client,
	}
}
