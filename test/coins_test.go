package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinsModel_UpdateCoinsData(t *testing.T) {
	err := coinsCtrl.Model.UpdateCoinsData(TestCoinData)
	assert.Nil(t, err)
}

func TestCoinsModel_GetCoinsData(t *testing.T) {
	coinsData, err := coinsCtrl.Model.GetCoinsData()
	assert.Nil(t, err)
	assert.NotZero(t, len(coinsData))
	assert.IsType(t, []models.Coin{}, coinsData)
	assert.Equal(t, TestCoinData, coinsData)
}

func TestCoinsController_GetCoinsAvailability(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	coins, err := coinsCtrl.GetCoinsAvailability(TestUser, c)
	assert.Nil(t, err)
	var coinsArray []models.Coin
	coinsBytes, err := json.Marshal(coins)
	assert.Nil(t, err)
	err = json.Unmarshal(coinsBytes, &coinsArray)
	assert.Equal(t, TestCoinData, coinsArray)
}

func TestCoinsController_UpdateCoinsAvailability(t *testing.T) {
	buf := new(bytes.Buffer)
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	token, err := utils.EncryptJWE(TestUser.ID, TestCoinData)
	assert.Nil(t, err)
	reqBody := models.BodyReq{
		Payload: token,
	}
	reqBytes, err := json.Marshal(reqBody)
	buf.Write(reqBytes)
	assert.Nil(t, err)
	_, err = resp.Write(reqBytes)
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest("POST", "/", buf)
	res, err := coinsCtrl.UpdateCoinsAvailability(TestUser, c)
	assert.Nil(t, err)
	assert.Equal(t, true, res)
}
