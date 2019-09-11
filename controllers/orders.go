package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
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
	userInfo, err := oc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []models.Order
	for _, id := range userInfo.Orders {
		obj, err := oc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (oc *OrdersController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("orderid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	userInfo, err := oc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Orders, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return oc.Model.Get(id)
}

func (oc *OrdersController) Store(userData models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (oc *OrdersController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	objs, err := oc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorAllError
	}
	return objs, nil
}

func (oc *OrdersController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("orderid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return oc.Model.Get(id)
}

func (oc *OrdersController) Update(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}
