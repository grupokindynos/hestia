package config

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)

var (
	ErrorNoAuth            = errors.New("you are not authorized")
	ErrorFbInitializeAuth  = errors.New("unable to initialize auth client")
	ErrorDbInitialize      = errors.New("unable to connect to database")
	ErrorNoUserInformation = errors.New("unable to get user information")
	ErrorMissingUID        = errors.New("missing user id param")
	ErrorUnmarshal         = errors.New("unable to unmarshal object")
	ErrorMissingID         = errors.New("missing id param")
	ErrorInfoDontMatchUser = errors.New("information requested doesn't match for this user")
	ErrorCoinDataGet       = errors.New("unable to get coin information")
	ErrorDecryptJWE        = errors.New("unable to decrypt jwe")
	ErrorDBStore           = errors.New("unable to store information to database")
	ErrorNotFound          = errors.New("information not found")
	ErrorAllError          = errors.New("something wrong happened, unable to get all records")
	ErrorAlreadyExists     = errors.New("object already exists")
	ErrorObol              = errors.New("unable to get obol rates")
	HttpClient             = &http.Client{
		Timeout: time.Second * 10,
	}
)

// GlobalResponseError is used to wrap all the errored API responses under the same model.
// Automatically detect if there is an error and return status and code according
func GlobalResponseError(result interface{}, err error, c *gin.Context) *gin.Context {
	if err != nil {
		c.JSON(500, gin.H{"message": "Error", "error": err.Error(), "status": -1})
	} else {
		c.JSON(200, gin.H{"data": result, "status": 1})
	}
	return c
}

// GlobalResponseNoAuth is used to wrap all non-auth API responses under the same model.
func GlobalResponseNoAuth(c *gin.Context) *gin.Context {
	c.JSON(401, gin.H{"message": "Error", "error": ErrorNoAuth.Error(), "status": -1})
	return c
}

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
