package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
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
	var configData hestia.Config
	coinsBytes, err := json.Marshal(coins)
	assert.Nil(t, err)
	err = json.Unmarshal(coinsBytes, &configData)
	assert.Equal(t, models.TestConfigData, configData)
}

func TestGlobalConfigController_UpdateConfigData(t *testing.T) {
	buf := new(bytes.Buffer)
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	token, err := jwt.EncryptJWE(models.TestUser.ID, models.TestConfigData)
	assert.Nil(t, err)
	reqBody := hestia.BodyReq{
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
