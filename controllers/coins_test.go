package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinsController_GetCoinsAvailability(t *testing.T) {
	coins, err := coinsCtrl.GetCoinsAvailability(models.TestUser, TestParams)
	assert.Nil(t, err)
	var coinsArray []hestia.Coin
	coinsBytes, err := json.Marshal(coins)
	assert.Nil(t, err)
	err = json.Unmarshal(coinsBytes, &coinsArray)
	assert.Equal(t, models.TestCoinData, coinsArray)
}

func TestCoinsController_UpdateCoinsAvailability(t *testing.T) {
	reqBytes, err := json.Marshal(models.TestCoinData)
	TestParams.Body = reqBytes
	res, err := coinsCtrl.UpdateCoinsAvailability(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.Equal(t, true, res)
}
