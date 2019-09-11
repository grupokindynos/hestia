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
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := CoinsModel{
		Db:         db,
		Collection: "coins",
	}
	coinsData, err := model.GetCoinsData()
	assert.Nil(t, err)
	assert.NotZero(t, len(coinsData))
	assert.IsType(t, []Coin{}, coinsData)
}

func TestCoinsModel_UpdateCoinsData(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := CoinsModel{
		Db:         db,
		Collection: "coins",
	}
	coinsData, err := model.GetCoinsData()
	assert.Nil(t, err)
	err = model.UpdateCoinsData(coinsData)
	assert.Nil(t, err)
}
