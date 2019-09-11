package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

var model CoinsModel

func init() {
	_ = godotenv.Load("../.env")
	db, _ := config.ConnectDB()
	model = CoinsModel{
		Db:         db,
		Collection: "coins",
	}
}

func TestCoinsModel_UpdateCoinsData(t *testing.T) {
	err := model.UpdateCoinsData(TestCoinData)
	assert.Nil(t, err)
}

func TestCoinsModel_GetCoinsData(t *testing.T) {
	coinsData, err := model.GetCoinsData()
	assert.Nil(t, err)
	assert.NotZero(t, len(coinsData))
	assert.IsType(t, []Coin{}, coinsData)
	assert.Equal(t, TestCoinData, coinsData)
}
