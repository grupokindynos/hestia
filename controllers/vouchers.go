package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
)

/*

	VouchersController is a safe-access query for vouchers on Firestore Database
	Database Structure:

	vouchers/
		VoucherID/
			voucherData

*/

type VouchersController struct {
	Model     *models.VouchersModel
	UserModel *models.UsersModel
}

// User methods

func (vc *VouchersController) GetUserAll(userData models.User, c *gin.Context) (interface{}, error) {
	userInfo, err := vc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []models.Voucher
	for _, id := range userInfo.Vouchers {
		obj, err := vc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (vc *VouchersController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("voucherid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	userInfo, err := vc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Orders, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return vc.Model.Get(id)
}

func (vc *VouchersController) Store(userData models.User, c *gin.Context) (interface{}, error) {
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
	var voucherData models.Voucher
	err = json.Unmarshal(rawBytes, &voucherData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	voucherData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(voucherData.PaymentData.Txid)))
	// Check if ID is already known on user data
	if utils.Contains(userData.Deposits, voucherData.ID) {
		return nil, config.ErrorAlreadyExists
	}
	// Check if ID is already known on data
	_, err = vc.Model.Get(voucherData.ID)
	if err == nil {
		return nil, config.ErrorAlreadyExists
	}
	voucherData.Status = "PENDING"
	err = vc.Model.Update(userData.ID, voucherData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	// Store ID on user information
	err = vc.UserModel.AddVoucher(userData.ID, voucherData.ID)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}

// Admin methods

func (vc *VouchersController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	objs, err := vc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorAllError
	}
	return objs, nil
}

func (vc *VouchersController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("voucherid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return vc.Model.Get(id)
}

func (vc *VouchersController) Update(userData models.User, c *gin.Context) (interface{}, error) {
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
	var voucherData models.Voucher
	err = json.Unmarshal(rawBytes, &voucherData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	// Hash the PaymentTxID as the ID
	// If this already exists, doesn't matter since it is deterministic
	voucherData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(voucherData.PaymentData.Txid)))
	// Store order data to process
	err = vc.Model.Update(userData.ID, voucherData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
