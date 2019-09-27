package config

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var (
	ErrorNoAuth            = errors.New("you are not authorized")
	ErrorFbInitializeAuth  = errors.New("unable to initialize auth client")
	ErrorDbInitialize      = errors.New("unable to connect to database")
	ErrorNoUserInformation = errors.New("unable to get user information")
	ErrorUnmarshal         = errors.New("unable to unmarshal object")
	ErrorMissingID         = errors.New("missing id param")
	ErrorInfoDontMatchUser = errors.New("information requested doesn't match for this user")
	ErrorCoinDataGet       = errors.New("unable to get coin information")
	ErrorConfigDataGet     = errors.New("unable to get config information")
	ErrorDecryptJWE        = errors.New("unable to decrypt jwe")
	ErrorDBStore           = errors.New("unable to store information to database")
	ErrorNotFound          = errors.New("information not found")
	ErrorAlreadyExists     = errors.New("object already exists")
)

// ConnectDB is the function to return the MongoDB connection.
func ConnectDB() (*mongo.Database, error) {
	MongoDB := os.Getenv("MONGODB_URL")
	MongoDBName := os.Getenv("MONGODB_NAME")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		return nil, err
	}
	db := client.Database(MongoDBName)
	return db, nil
}
