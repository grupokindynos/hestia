package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/models"
	"os"
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

func (oc *OrdersController) GetAll(userData hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		return oc.Model.GetAll(params.Filter)
	}
	userInfo, err := oc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.Order
	for _, id := range userInfo.Orders {
		obj, err := oc.Model.Get(id)
		if err != nil {
			return nil, errors.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (oc *OrdersController) GetSingle(userData hestia.User, params Params) (interface{}, error) {
	if params.OrderID == "" {
		return nil, errors.ErrorMissingID
	}
	if params.Admin {
		return oc.Model.Get(params.OrderID)
	}
	userInfo, err := oc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Orders, params.OrderID) {
		return nil, errors.ErrorInfoDontMatchUser
	}
	return oc.Model.Get(params.OrderID)
}

func (oc *OrdersController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody hestia.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Orders Microservice signature
	rawBytes, err := jwt.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var orderData hestia.Order
	err = json.Unmarshal(rawBytes, &orderData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	orderData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(orderData.PaymentInfo.Txid)))
	// Check if ID is already known on data
	_, err = oc.Model.Get(orderData.ID)
	if err == nil {
		responses.GlobalResponseError(nil, errors.ErrorAlreadyExists, c)
		return
	}
	orderData.Status = "PENDING"
	err = oc.Model.Update(orderData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = oc.UserModel.AddVoucher(orderData.UID, orderData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	response, err := jwt.EncodeJWS(orderData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseError(response, err, c)
	return
}
