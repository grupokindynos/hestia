package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestOrdersController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	orders, err := ordersCtrl.GetUserAll(models.TestUser, c)
	assert.Nil(t, err)
	var ordersArray []models.Order
	orderBytes, err := json.Marshal(orders)
	assert.Nil(t, err)
	err = json.Unmarshal(orderBytes, &ordersArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Order{}, orders)
	assert.Equal(t, models.TestOrder, ordersArray[0])
}

func TestOrdersController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "orderid", Value: models.TestOrder.ID}}
	order, err := ordersCtrl.GetUserSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Order{}, order)
	assert.Equal(t, models.TestOrder, order)
}

func TestOrdersController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	orders, err := ordersCtrl.GetAll(models.TestUser, c)
	assert.Nil(t, err)
	var orderArray []models.Order
	orderBytes, err := json.Marshal(orders)
	assert.Nil(t, err)
	err = json.Unmarshal(orderBytes, &orderArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Order{}, orders)
	assert.Equal(t, models.TestOrder, orderArray[0])
}

func TestOrdersController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "orderid", Value: models.TestOrder.ID}}
	order, err := ordersCtrl.GetSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Order{}, order)
	assert.Equal(t, models.TestOrder, order)
}
