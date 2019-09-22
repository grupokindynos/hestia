package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/jws"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGlobalConfigController_GetConfig(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	coins, err := globalCtrl.GetConfig(models.TestUser, c, false)
	assert.Nil(t, err)
	var coinsArray []models.Coin
	coinsBytes, err := json.Marshal(coins)
	assert.Nil(t, err)
	err = json.Unmarshal(coinsBytes, &coinsArray)
	assert.Equal(t, models.TestCoinData, coinsArray)
}

func TestGlobalConfigController_UpdateConfigData(t *testing.T) {
	buf := new(bytes.Buffer)
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	token, err := jws.EncryptJWE(models.TestUser.ID, models.TestConfigData)
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
	res, err := globalCtrl.UpdateConfigData(models.TestUser, c, false)
	assert.Nil(t, err)
	assert.Equal(t, true, res)
}
