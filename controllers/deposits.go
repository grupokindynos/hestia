package controllers

import (
	"encoding/json"
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

func (dc *DepositsController) GetAll(userData hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		return dc.Model.GetAll(params.Filter)
	}
	userInfo, err := dc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.Deposit
	for _, id := range userInfo.Deposits {
		obj, err := dc.Model.Get(id)
		if err != nil {
			return nil, errors.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (dc *DepositsController) GetSingle(userData hestia.User, params Params) (interface{}, error) {
	if params.DepositID == "" {
		return nil, errors.ErrorMissingID
	}
	if params.Admin {
		return dc.Model.Get(params.DepositID)
	}
	userInfo, err := dc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Deposits, params.DepositID) {
		return nil, errors.ErrorInfoDontMatchUser
	}
	return dc.Model.Get(params.DepositID)
}

func (dc *DepositsController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody hestia.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Deposits Microservice signature
	rawBytes, err := jwt.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var depositData hestia.Deposit
	err = json.Unmarshal(rawBytes, &depositData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	err = dc.Model.Update(depositData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	response, err := jwt.EncodeJWS(depositData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseError(response, err, c)
	return
}
