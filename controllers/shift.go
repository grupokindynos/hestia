package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
	"os"
)

/*

	ShiftController is a safe-access query for shifts on Firestore Database
	Database Structure:

	shift/
		ShiftID/
			shiftData

*/

type ShiftsController struct {
	Model     *models.ShiftModel
	UserModel *models.UsersModel
}

func (sc *ShiftsController) GetAll(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return sc.Model.GetAll()
	}
	userInfo, err := sc.UserModel.Get(userData.ID)
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

func (sc *ShiftsController) GetSingle(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	id, ok := c.Params.Get("shiftid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	if admin {
		return sc.Model.Get(id)
	}
	userInfo, err := sc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Shifts, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return sc.Model.Get(id)
}

func (sc *ShiftsController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Try to decrypt it
	rawBytes, err := utils.DecodeJWS(ReqBody.Payload, os.Getenv("TYCHE_PUBLIC_KEY"))
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var shiftData models.Shift
	err = json.Unmarshal(rawBytes, &shiftData)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	shiftData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(shiftData.Payment.Txid)))
	// Check if ID is already known on data
	_, err = sc.Model.Get(shiftData.ID)
	if err == nil {
		config.GlobalResponseError(nil, config.ErrorAlreadyExists, c)
		return
	}
	// Store shift data to process
	err = sc.Model.Update(shiftData)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = sc.UserModel.AddShift(shiftData.UID, shiftData.ID)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	response, err := utils.EncodeJWS(shiftData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	config.GlobalResponseError(response, err, c)
	return
}
