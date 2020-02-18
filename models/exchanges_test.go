package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

/* func TestExchangesModel_Update(t *testing.T) {
	err := exchangesModel.Update(TestExchangeData)
	assert.Nil(t, err)
} */

func TestExchangesModel_GetAll(t *testing.T) {
	orders, err := exchangesModel.GetAll(true, 0)
	assert.Nil(t, err)
	assert.NotZero(t, len(orders))
	// assert.Equal(t, TestExchangeData, orders[0])
	assert.IsType(t, []hestia.AdrestiaOrder{}, orders)
}

func TestExchangesModel_Get(t *testing.T) {
	order, err := exchangesModel.Get(TestExchangeData.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestExchangeData, order)
}
