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

func NewMongodbClient(ctx context.Context, config *utils.Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@db:%s", config.DbUsername, config.DbPassword, config.DbPort)
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}

func SetupMongodb(client *mongo.Client) error {
	coll := client.Database(DB).Collection(UsersCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", -1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}
