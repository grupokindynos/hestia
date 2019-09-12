package config

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	_ = godotenv.Load("../.env")
}
func TestGlobalResponseError(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	newErr := errors.New("test error")
	_ = GlobalResponseError(nil, newErr, c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.Nil(t, response["data"])
	assert.Equal(t, newErr.Error(), response["error"])
	assert.Equal(t, float64(-1), response["status"])
}

func TestGlobalResponseError2(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	newErr := "test success"
	_ = GlobalResponseError(newErr, nil, c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Nil(t, response["error"])
	assert.Equal(t, newErr, response["data"])
	assert.Equal(t, float64(1), response["status"])
}

func TestGlobalResponseNoAuth(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	_ = GlobalResponseNoAuth(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.Code)
	assert.Nil(t, response["data"])
	assert.Equal(t, ErrorNoAuth.Error(), response["error"])
	assert.Equal(t, float64(-1), response["status"])
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
