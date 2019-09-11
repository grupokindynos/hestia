package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestOrdersModel_Get(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	newOrder, err := model.Get(TestOrder.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestOrder, newOrder)
}

func TestOrdersModel_GetAll(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	orders, err := model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(orders))
	assert.IsType(t, []Order{}, orders)
}
