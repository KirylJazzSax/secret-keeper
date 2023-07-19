package db

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB                = "secret-keeper"
	UsersCollection   = "users"
	SecretsCollection = "secrets"
)

type Db struct {
	Client *mongo.Client
}

func (db *Db) Shutdown() error {
	return db.Client.Disconnect(context.TODO())
}

func NewDbClient(ctx context.Context, config *utils.Config) (*Db, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@db:%s", config.DbUsername, config.DbPassword, config.DbPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	return &Db{
		Client: client,
	}, nil
}

func SetupMongodb(ctx context.Context, client *mongo.Client) error {
	coll := client.Database(DB).Collection(UsersCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", -1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	return err
}
