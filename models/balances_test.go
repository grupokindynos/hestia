package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalancesModel_UpdateBalances(t *testing.T) {
	err := balancesModel.UpdateBalances(TestBalances)
	assert.Nil(t, err)
}

func TestBalancesModel_GetBalances(t *testing.T) {
	balanceData, err := balancesModel.GetBalances()
	assert.Nil(t, err)
	// assert.NotZero(t, len(balanceData))
	assert.IsType(t, []hestia.CoinBalances{}, balanceData)
	// assert.Equal(t, TestBalances, balanceData)
}
