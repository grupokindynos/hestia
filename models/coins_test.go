package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)


func init() {
	_ = godotenv.Load("../.env")
}

func TestCoinsModel_GetCoinsData(t *testing.T) {
	db, _ := config.ConnectDB()
	model := CoinsModel{
		Db:         db,
		Collection: "coins",
	}
	coinsData, err := model.GetCoinsData()
	assert.Nil(t, err)
	assert.NotZero(t, len(coinsData))
	assert.IsType(t, []Coin{}, coinsData)
	assert.Equal(t, TestCoinData, coinsData)
}
