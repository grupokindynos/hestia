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

	DepositsController is a safe-access query for deposits on Firestore Database
	Database Structure:

	deposits/
		DepositID/
			depositData

*/

type DepositsController struct {
	Model     *models.DepositsModel
	UserModel *models.UsersModel
}

func (dc *DepositsController) GetAll(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return dc.Model.GetAll()
	}
	userInfo, err := dc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []hestia.Deposit
	for _, id := range userInfo.Deposits {
		obj, err := dc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (dc *DepositsController) GetSingle(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	id, ok := c.Params.Get("depositid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	if admin {
		return dc.Model.Get(id)
	}
	userInfo, err := dc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Deposits, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return dc.Model.Get(id)
}

func (dc *DepositsController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody hestia.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Deposits Microservice signature
	rawBytes, err := jwt.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var depositData hestia.Deposit
	err = json.Unmarshal(rawBytes, &depositData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	// If this already exists, doesn't matter since it is deterministic
	depositData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(depositData.Payment.Txid)))
	// Store deposit data to process
	err = dc.Model.Update(depositData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	response, err := jwt.EncodeJWS(depositData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseError(response, err, c)
	return
}
