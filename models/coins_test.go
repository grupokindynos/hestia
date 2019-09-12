package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinsModel_UpdateCoinsData(t *testing.T) {
	err := coinsModel.UpdateCoinsData(TestCoinData)
	assert.Nil(t, err)
}

func TestCoinsModel_GetCoinsData(t *testing.T) {
	coinsData, err := coinsModel.GetCoinsData()
	assert.Nil(t, err)
	assert.NotZero(t, len(coinsData))
	assert.IsType(t, []Coin{}, coinsData)
	assert.Equal(t, TestCoinData, coinsData)
}
