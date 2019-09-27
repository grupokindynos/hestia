package config

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

func init() {
	_ = godotenv.Load("../.env")
}

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.IsType(t, &mongo.Database{}, db)
}

func TestConnectDBError(t *testing.T) {
	err := os.Setenv("MONGODB_URL", "")
	assert.Nil(t, err)
	db, err := ConnectDB()
	assert.Nil(t, db)
	assert.NotNil(t, err)
}
