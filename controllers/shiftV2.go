package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/models"
	"os"
	"strconv"
)

/*

	ShiftController is a safe-access query for shifts on Firestore Database
	Database Structure:

	shift/
		ShiftID/
			shiftData

*/

type ShiftsControllerV2 struct {
	Model     *models.ShiftModelV2
	UserModel *models.UsersModel
}

func (sc *ShiftsControllerV2) GetAll(userData hestia.User, params Params) (interface{}, error) {
	filterNum, _ := strconv.ParseInt(params.Filter, 10, 32)
	if params.Filter == "all" {
		filterNum = -1
	}
	if params.Admin {
		return sc.Model.GetAll(int32(filterNum), "")
	}
	userInfo, err := sc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.LightShift
	for _, id := range userInfo.ShiftV2 {
		obj, err := sc.Model.Get(id)
		if err != nil {
			return nil, errors.ErrorNotFound
		}
		var newShift = hestia.LightShift{
			ID:        obj.ID,
			UID:       obj.UID,
			Status:    hestia.GetShiftStatusv2String(obj.Status),
			Timestamp: obj.Timestamp,
			Payment: hestia.LightPayment{
				Address: obj.Payment.Address,
				Coin:    obj.Payment.Coin,
				Txid:    obj.Payment.Txid,
				Amount:  obj.Payment.Amount,
			},
			RefundAddr:         obj.RefundAddr,
			ToCoin:             obj.ToCoin,
			ToAmount:           obj.ToAmount,
			UserReceivedAmount: obj.UserReceivedAmount,
			ToAddress:          obj.ToAddress,
			PaymentProof:       obj.PaymentProof,
			ProofTimestamp:     obj.ProofTimestamp,
		}
		fmt.Println(newShift)
		Array = append(Array, newShift)
	}
	return Array, nil
}

func (sc *ShiftsControllerV2) GetSingle(userData hestia.User, params Params) (interface{}, error) {
	if params.ShiftID == "" {
		return nil, errors.ErrorMissingID
	}
	if params.Admin {
		return sc.Model.Get(params.ShiftID)
	}
	userInfo, err := sc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Shifts, params.ShiftID) {
		return nil, errors.ErrorInfoDontMatchUser
	}
	return sc.Model.Get(params.ShiftID)
}

func (sc *ShiftsControllerV2) GetSingleTyche(c *gin.Context) {
	// Check if the user has an id
	id, ok := c.Params.Get("shiftid")
	if !ok {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	_, err := mvt.VerifyRequest(c)
	if err != nil {
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

func (sc *ShiftsControllerV2) GetAllTyche(c *gin.Context) {
	filter := c.Query("filter")
	filterNum, _ := strconv.ParseInt(filter, 10, 32)
	if filter == "" {
		filterNum = -1
	}

	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	shiftList, err := sc.Model.GetAll(int32(filterNum), "")
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), shiftList, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (sc *ShiftsControllerV2) Store(c *gin.Context) {
	payload, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var shiftData hestia.ShiftV2
	err = json.Unmarshal(payload, &shiftData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Store shift data to process
	err = sc.Model.Update(shiftData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = sc.UserModel.AddShiftV2(shiftData.UID, shiftData.ID)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), shiftData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}
