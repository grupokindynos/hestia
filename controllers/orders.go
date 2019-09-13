package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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

func (oc *OrdersController) GetAll(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return oc.Model.GetAll()
	}
	userInfo, err := oc.UserModel.Get(userData.ID)
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

func (oc *OrdersController) GetSingle(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	id, ok := c.Params.Get("orderid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	if admin {
		return oc.Model.Get(id)
	}
	userInfo, err := oc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Orders, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return oc.Model.Get(id)
}

func (oc *OrdersController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Try to decrypt it
	rawBytes, err := utils.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, config.ErrorDecryptJWE
	}
	// Try to unmarshal the information of the payload
	var orderData models.Order
	err = json.Unmarshal(rawBytes, &orderData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	orderData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(orderData.PaymentInfo.Txid)))
	// Check if ID is already known on user data
	if utils.Contains(userData.Deposits, orderData.ID) {
		return nil, config.ErrorAlreadyExists
	}
	// Check if ID is already known on data
	_, err = oc.Model.Get(orderData.ID)
	if err == nil {
		return nil, config.ErrorAlreadyExists
	}
	orderData.Status = "PENDING"
	err = oc.Model.Update(orderData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	// Store ID on user information
	err = oc.UserModel.AddVoucher(userData.ID, orderData.ID)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
