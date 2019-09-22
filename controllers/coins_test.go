package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jws"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinsController_GetCoinsAvailability(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	coins, err := coinsCtrl.GetCoinsAvailability(models.TestUser, c, false)
	assert.Nil(t, err)
	var coinsArray []hestia.Coin
	coinsBytes, err := json.Marshal(coins)
	assert.Nil(t, err)
	err = json.Unmarshal(coinsBytes, &coinsArray)
	assert.Equal(t, models.TestCoinData, coinsArray)
}

func TestCoinsController_UpdateCoinsAvailability(t *testing.T) {
	buf := new(bytes.Buffer)
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	token, err := jws.EncryptJWE(models.TestUser.ID, models.TestCoinData)
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
	res, err := coinsCtrl.UpdateCoinsAvailability(models.TestUser, c, false)
	assert.Nil(t, err)
	assert.Equal(t, true, res)
}
