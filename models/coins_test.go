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

var TestCoinData = []Coin{
	{"BTC", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
	{"LTC", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
	{"DASH", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
	{"POLIS", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
	{"GRS", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
	{"DGB", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
	{"COLX", false, false, false, false, Balances{HotWallet: 1, Exchanges: 1}},
}

func TestCoinsModel_UpdateCoinsData(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := CoinsModel{
		Db:         db,
		Collection: "coins",
	}
	assert.Nil(t, err)
	err = model.UpdateCoinsData(TestCoinData)
	assert.Nil(t, err)
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
	assert.Equal(t, TestCoinData, coinsData)
}
