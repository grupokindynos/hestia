package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestOrdersModel_Update(t *testing.T) {
	err := ordersCtrl.Model.Update(TestOrder)
	assert.Nil(t, err)
}

func TestOrdersModel_Get(t *testing.T) {
	newOrder, err := ordersCtrl.Model.Get(TestOrder.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestOrder, newOrder)
}

func TestOrdersModel_GetAll(t *testing.T) {
	orders, err := ordersCtrl.Model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(orders))
	assert.IsType(t, []models.Order{}, orders)
}

func TestOrdersController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	orders, err := ordersCtrl.GetUserAll(TestUser, c)
	assert.Nil(t, err)
	var ordersArray []models.Order
	orderBytes, err := json.Marshal(orders)
	assert.Nil(t, err)
	err = json.Unmarshal(orderBytes, &ordersArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Order{}, orders)
	assert.Equal(t, TestOrder, ordersArray[0])
}

func TestOrdersController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "orderid", Value: TestOrder.ID}}
	order, err := ordersCtrl.GetUserSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Order{}, order)
	assert.Equal(t, TestOrder, order)
}

func TestOrdersController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	orders, err := ordersCtrl.GetAll(TestUser, c)
	assert.Nil(t, err)
	var orderArray []models.Order
	orderBytes, err := json.Marshal(orders)
	assert.Nil(t, err)
	err = json.Unmarshal(orderBytes, &orderArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Order{}, orders)
	assert.Equal(t, TestOrder, orderArray[0])
}

func TestOrdersController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "orderid", Value: TestOrder.ID}}
	order, err := ordersCtrl.GetSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Order{}, order)
	assert.Equal(t, TestOrder, order)
}
