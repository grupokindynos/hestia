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
	"strconv"
)

/*

	ShiftController is a safe-access query for shifts on Firestore Database
	Database Structure:

	shift/
		ShiftID/
			shiftData

*/

type ShiftsController struct {
	Obol      *services.ObolService
	Model     *models.ShiftModel
	UserModel *models.UsersModel
}

// User methods

func (sc *ShiftsController) GetUserAll(userData models.User, c *gin.Context) (interface{}, error) {
	userInfo, err := sc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []models.Shift
	for _, id := range userInfo.Shifts {
		obj, err := sc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (sc *ShiftsController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("shiftid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	userInfo, err := sc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Shifts, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return sc.Model.Get(id)
}

func (sc *ShiftsController) Store(userData models.User, c *gin.Context) (interface{}, error) {
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
	var shiftData models.Shift
	err = json.Unmarshal(rawBytes, &shiftData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	shiftData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(shiftData.Payment.Txid)))
	// Check if ID is already known on user data
	if utils.Contains(userData.Shifts, shiftData.ID) {
		return nil, config.ErrorAlreadyExists
	}
	// Check if iD is already known on data
	_, err = sc.Model.Get(shiftData.ID)
	if err == nil {
		return nil, config.ErrorAlreadyExists
	}

	conversionRate, err := sc.Obol.GetComplexRate(shiftData.Payment.Coin, shiftData.Conversion.Coin, shiftData.Payment.Amount)
	if err != nil {
		return nil, config.ErrorObol
	}
	sendAmount, _ := strconv.ParseFloat(shiftData.Payment.Amount, 64)
	// Fill missing data
	shiftData.Payment = models.Payment{
		Address:       shiftData.Conversion.Address,
		Amount:        fmt.Sprintf("%f", conversionRate*sendAmount),
		Coin:          shiftData.Conversion.Coin,
		Txid:          "",
		Confirmations: "",
	}
	shiftData.Status = "PENDING"
	// Store shift data to process
	err = sc.Model.Update(shiftData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	// Store ID on user information
	err = sc.UserModel.AddShift(userData.ID, shiftData.ID)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}

// Admin methods

func (sc *ShiftsController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	objs, err := sc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorAllError
	}
	return objs, nil
}

func (sc *ShiftsController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("shiftid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return sc.Model.Get(id)
}

func (sc *ShiftsController) Update(userData models.User, c *gin.Context) (interface{}, error) {
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
	var shiftData models.Shift
	err = json.Unmarshal(rawBytes, &shiftData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	// If this already exists, doesn't matter since it is deterministic
	shiftData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(shiftData.Payment.Txid)))
	// Store shift data to process
	err = sc.Model.Update(shiftData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
