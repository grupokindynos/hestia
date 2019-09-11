package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var orderCtrl OrdersController

func init() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	model := &models.OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	userModel := &models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	orderCtrl = OrdersController{
		Model:     model,
		UserModel: userModel,
	}
}

func TestOrdersController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	orders, err := orderCtrl.GetUserAll(models.TestUser, c)
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
	order, err := orderCtrl.GetUserSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Order{}, order)
	assert.Equal(t, models.TestOrder, order)
}

func TestOrdersController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	orders, err := orderCtrl.GetAll(models.TestUser, c)
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
	order, err := orderCtrl.GetSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Order{}, order)
	assert.Equal(t, models.TestOrder, order)
}
