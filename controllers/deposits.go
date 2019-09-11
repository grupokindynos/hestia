package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services"
	"github.com/grupokindynos/hestia/utils"
)

/*

	DepositsController is a safe-access query for deposits on Firestore Database
	Database Structure:

	deposits/
		DepositID/
			depositData

*/

type DepositsController struct {
	Obol      *services.ObolService
	Model     *models.DepositsModel
	UserModel *models.UsersModel
}

// User methods

func (dc *DepositsController) GetUserAll(userData models.User, c *gin.Context) (interface{}, error) {
	userInfo, err := dc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []models.Deposit
	for _, id := range userInfo.Deposits {
		obj, err := dc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (dc *DepositsController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("depositid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	userInfo, err := dc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Deposits, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return dc.Model.Get(id)
}

func (dc *DepositsController) Store(userData models.User, c *gin.Context) (interface{}, error) {
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
	var depositData models.Deposit
	err = json.Unmarshal(rawBytes, &depositData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	depositData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(depositData.Payment.Txid)))
	// Check if ID is already known on user data
	if utils.Contains(userData.Deposits, depositData.ID) {
		return nil, config.ErrorAlreadyExists
	}
	// Check if ID is already known on data
	_, err = dc.Model.Get(depositData.ID)
	if err == nil {
		return nil, config.ErrorAlreadyExists
	}
	rate, err := dc.Obol.GetSimpleRate(depositData.Payment.Coin)
	if err != nil {
		return nil, config.ErrorObol
	}
	depositData.AmountInPeso = fmt.Sprintf("%f", rate)
	depositData.Status = "PENDING"
	err = dc.Model.Update(userData.ID, depositData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	// Store ID on user information
	err = dc.UserModel.AddDeposit(userData.ID, depositData.ID)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}

// Admin methods

func (dc *DepositsController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	objs, err := dc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorAllError
	}
	return objs, nil
}

func (dc *DepositsController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("depositid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return dc.Model.Get(id)
}

func (dc *DepositsController) Update(userData models.User, c *gin.Context) (interface{}, error) {
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
	var depositData models.Deposit
	err = json.Unmarshal(rawBytes, &depositData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	// If this already exists, doesn't matter since it is deterministic
	depositData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(depositData.Payment.Txid)))
	// Store deposit data to process
	err = dc.Model.Update(userData.ID, depositData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
