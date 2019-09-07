package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	OrdersController is a safe-access query for orders on Firestore Database
	Database Structure:

	orders/
		UID/
          	orders -> Array of OrderIDs

	orderIndex/
		OrderID/
			orderData

*/

type OrdersController struct {
	Model *models.OrdersModel
}

// User methods

func (oc *OrdersController) GetUserAll(uid string, params gin.Params) (interface{}, error) {
	return nil, nil
}

func (oc *OrdersController) GetUserSingle(uid string, params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}

func (oc *OrdersController) Store(uid string, params gin.Params) (interface{}, error) {
	return "", nil
}

// Admin methods

func (oc *OrdersController) GetAll(params gin.Params) (interface{}, error) {
	return nil, nil
}

func (oc *OrdersController) GetSingle(params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}
