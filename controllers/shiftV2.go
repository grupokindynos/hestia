package controllers

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/models"
)

/*

	ShiftController is a safe-access query for shifts on Firestore Database
	Database Structure:

	shifts2/
		ShiftID/
			shiftData

*/

type ShiftsControllerV2 struct {
	Model       *models.ShiftModelV2
	LegacyModel *models.ShiftModel
	UserModel   *models.UsersModel
}

func (sc *ShiftsControllerV2) GetAll(userData hestia.User, params Params) (interface{}, error) {
	log.Println("retrieving shift data for ", userData.ID)
	filterNum, _ := strconv.ParseInt(params.Filter, 10, 32)
	if params.Filter == "all" {
		filterNum = -1
	}
	if params.Admin {
		return sc.Model.GetAll(int32(filterNum), "")
	}
	userInfo, err := sc.UserModel.Get(userData.ID)
	if err != nil {
		log.Println("ShiftV2::GetAll::NOUSERINFO::", userData.ID)
		return nil, errors.ErrorNoUserInformation
	}
	timestamp, err := strconv.ParseInt(params.Timestamp, 10, 64)
	if err != nil {
		timestamp = 0
	}

	var lastTimestamp int64
	var Array []hestia.LightShift
	// shifts from v2
	for _, id := range userInfo.ShiftV2 {
		obj, err := sc.Model.Get(id)
		if err != nil {
			continue
		}
		if obj.LastUpdated == 0 {
			_ = sc.Model.Update(obj)
		} else if obj.LastUpdated < timestamp {continue}

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
			ToAmount:           int64(obj.UserReceivedAmount * 10e8),
			UserReceivedAmount: obj.UserReceivedAmount,
			ToAddress:          obj.ToAddress,
			PaymentProof:       obj.PaymentProof,
			ProofTimestamp:     obj.ProofTimestamp,
		}
		Array = append(Array, newShift)
	}

	lastTimestamp = time.Now().Unix()
	// shifts from v1
	for _, id := range userInfo.Shifts {
		obj, err := sc.LegacyModel.Get(id)
		if err != nil {
			continue
		}
		if obj.Timestamp < timestamp {continue}
		var newShift = hestia.LightShift{
			ID:        obj.ID,
			UID:       obj.UID,
			Status:    obj.Status,
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
			UserReceivedAmount: float64(obj.ToAmount) * 1e-8,
			ToAddress:          obj.ToAddress,
			PaymentProof:       obj.PaymentProof,
			ProofTimestamp:     obj.ProofTimestamp,
		}
		Array = append(Array, newShift)
	}
	if timestamp == 0 {
		return Array, nil
	} else {
		return hestia.ShiftHistoryResponse{
			Shifts:    Array,
			Timestamp: lastTimestamp,
		}, nil
	}
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
	_, _, err := mvt.VerifyRequest(c)
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

	_, _, err := mvt.VerifyRequest(c)
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

func (sc *ShiftsControllerV2) GetOpenShifts(c *gin.Context) {
	timestamp, err := strconv.ParseInt(c.Query("timestamp"), 10, 64)
	if err != nil {
		ysd := time.Now()
		ysd.AddDate(0, 0, -1)
		timestamp = time.Now().AddDate(0, 0, -1).Unix()
	}

	_, _, err = mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	shifts, err := sc.Model.GetOpenShifts(timestamp)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), shifts, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (sc *ShiftsControllerV2) Store(c *gin.Context) {
	payload, _, err := mvt.VerifyRequest(c)
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

func (sc *ShiftsControllerV2) GetShiftsByTimestampTyche(c *gin.Context) {
	userId := c.Query("userid")
	if userId == "" {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	ts := c.Query("timestamp")
	if ts == "" {
		responses.GlobalResponseError(nil, errors.ErrorMissingID, c)
		return
	}
	_, _, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	userInfo, err := sc.UserModel.Get(userId)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorNoUserInformation, c)
		return
	}
	var userShifts []hestia.ShiftV2
	timestamp, _ := strconv.ParseInt(ts, 10, 64)

	for _, id := range userInfo.ShiftV2 {
		obj, err := sc.Model.Get(id)
		if err != nil {
			continue
		}

		if timestamp <= obj.Timestamp {
			userShifts = append(userShifts, obj)
		}
	}

	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), userShifts, os.Getenv("HESTIA_PRIVATE_KEY"))

	responses.GlobalResponseMRT(header, body, c)
	return

}
