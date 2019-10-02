package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrdersController_GetUserAll(t *testing.T) {
	orders, err := ordersCtrl.GetAll(models.TestUser, TestParams)
	assert.Nil(t, err)
	var ordersArray []hestia.Order
	orderBytes, err := json.Marshal(orders)
	assert.Nil(t, err)
	err = json.Unmarshal(orderBytes, &ordersArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Order{}, orders)
	assert.Equal(t, models.TestOrder, ordersArray[0])
}

func TestOrdersController_GetUserSingle(t *testing.T) {
	order, err := ordersCtrl.GetSingle(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Order{}, order)
	assert.Equal(t, models.TestOrder, order)
}

func TestOrdersController_GetAll(t *testing.T) {
	orders, err := ordersCtrl.GetAll(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	var orderArray []hestia.Order
	orderBytes, err := json.Marshal(orders)
	assert.Nil(t, err)
	err = json.Unmarshal(orderBytes, &orderArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Order{}, orders)
	assert.Equal(t, models.TestOrder, orderArray[0])
}

func TestOrdersController_GetSingle(t *testing.T) {
	order, err := ordersCtrl.GetSingle(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Order{}, order)
	assert.Equal(t, models.TestOrder, order)
}
