package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/config"
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

func (oc *OrdersController) GetAll(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return oc.Model.GetAll()
	}
	userInfo, err := oc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []hestia.Order
	for _, id := range userInfo.Orders {
		obj, err := oc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (oc *OrdersController) GetSingle(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
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
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Orders Microservice signature
	rawBytes, err := jwt.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var orderData hestia.Order
	err = json.Unmarshal(rawBytes, &orderData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	orderData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(orderData.PaymentInfo.Txid)))
	// Check if ID is already known on data
	_, err = oc.Model.Get(orderData.ID)
	if err == nil {
		responses.GlobalResponseError(nil, config.ErrorAlreadyExists, c)
		return
	}
	orderData.Status = "PENDING"
	err = oc.Model.Update(orderData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = oc.UserModel.AddVoucher(orderData.UID, orderData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	response, err := jwt.EncodeJWS(orderData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseError(response, err, c)
	return
}
