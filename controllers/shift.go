package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
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

func (sc *ShiftsController) GetAll(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return sc.Model.GetAll("all")
	}
	userInfo, err := sc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []hestia.Shift
	for _, id := range userInfo.Shifts {
		obj, err := sc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (sc *ShiftsController) GetSingle(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
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

func (sc *ShiftsController) GetSingleTyche(c *gin.Context) {
	// Check if the user has an id
	id, ok := c.Params.Get("shiftid")
	if !ok {
		responses.GlobalResponseError(nil, config.ErrorMissingID, c)
		return
	}
	headerSignature := os.Getenv("service")
	if headerSignature == "" {
		responses.GlobalResponseNoAuth(c)
		return
	}
	valid, _ := mvt.VerifyMVTToken(headerSignature, nil, os.Getenv("TYCHE_PUBLIC_KEY"), os.Getenv("MASTER_PASSWORD"))
	if !valid {
		responses.GlobalResponseNoAuth(c)
		return
	}
	shift, err := sc.Model.Get(id)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), shift, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (sc *ShiftsController) GetAllTyche(c *gin.Context) {
	filter := c.Query("filter")
	if filter == "" {
		filter = "all"
	}
	headerSignature := os.Getenv("service")
	if headerSignature == "" {
		responses.GlobalResponseNoAuth(c)
		return
	}
	valid, _ := mvt.VerifyMVTToken(headerSignature, nil, os.Getenv("TYCHE_PUBLIC_KEY"), os.Getenv("MASTER_PASSWORD"))
	if !valid {
		responses.GlobalResponseNoAuth(c)
		return
	}
	shiftList, err := sc.Model.GetAll(filter)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), shiftList, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (sc *ShiftsController) Store(c *gin.Context) {
	// Catch the request body
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	headerSignature := os.Getenv("service")
	if headerSignature == "" {
		responses.GlobalResponseNoAuth(c)
		return
	}
	reqBytes, err := json.Marshal(ReqBody.Payload)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	valid, payload := mvt.VerifyMVTToken(headerSignature, reqBytes, os.Getenv("TYCHE_PUBLIC_KEY"), os.Getenv("MASTER_PASSWORD"))
	if !valid {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var shiftData hestia.Shift
	err = json.Unmarshal(payload, &shiftData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	shiftData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(shiftData.Payment.Txid)))
	// Check if ID is already known on data
	_, err = sc.Model.Get(shiftData.ID)
	if err == nil {
		responses.GlobalResponseError(nil, config.ErrorAlreadyExists, c)
		return
	}
	// Store shift data to process
	err = sc.Model.Update(shiftData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = sc.UserModel.AddShift(shiftData.UID, shiftData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), shiftData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}
