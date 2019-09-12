package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrdersModel_Update(t *testing.T) {
	err := ordersModel.Update(TestOrder)
	assert.Nil(t, err)
}

func TestOrdersModel_Get(t *testing.T) {
	newOrder, err := ordersModel.Get(TestOrder.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestOrder, newOrder)
}

func TestOrdersModel_GetAll(t *testing.T) {
	orders, err := ordersModel.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(orders))
	assert.IsType(t, []Order{}, orders)
}
