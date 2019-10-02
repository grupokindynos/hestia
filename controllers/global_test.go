package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobalConfigController_GetConfig(t *testing.T) {
	coins, err := globalCtrl.GetConfig(models.TestUser, TestParams)
	assert.Nil(t, err)
	var configData hestia.Config
	coinsBytes, err := json.Marshal(coins)
	assert.Nil(t, err)
	err = json.Unmarshal(coinsBytes, &configData)
	assert.Equal(t, models.TestConfigData, configData)
}

func TestGlobalConfigController_UpdateConfigData(t *testing.T) {
	reqBytes, err := json.Marshal(models.TestConfigData)
	TestParams.Body = reqBytes
	res, err := globalCtrl.UpdateConfigData(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.Equal(t, true, res)
}
