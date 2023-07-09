package repository

import (
	"context"
	"math/big"
	"time"

	"github.com/KirylJazzSax/secret-keeper/internal/common/db"
	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDto struct {
	Id        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	Password  string             `bson:"password"`
}

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
	user := &UserDto{}

	err := coll.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	n := new(big.Int)
	id, err := n.SetString(user.Id.Hex(), 16)

	u := &domain.User{
		Id:        id.Int64(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return u, err
}

func NewMongoUserRepository(client *mongo.Client) *MongoUserRepository {
	return &MongoUserRepository{
		client: client,
	}
}
