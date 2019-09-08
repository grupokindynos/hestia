package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	OrdersController is a safe-access query for orders on Firestore Database
	Database Structure:

	orders/
		OrderID/
			orderData

*/

type OrdersController struct {
	Model     *models.OrdersModel
	UserModel *models.UsersModel
}

// User methods

func (oc *OrdersController) GetUserAll(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (oc *OrdersController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}

func (oc *OrdersController) Store(userData models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (oc *OrdersController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (oc *OrdersController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}
