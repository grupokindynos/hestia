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

func (vc *VouchersController) GetAll(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return vc.Model.GetAll()
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
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

func (vc *VouchersController) GetSingle(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	id, ok := c.Params.Get("voucherid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	if admin {
		return vc.Model.Get(id)
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Orders, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return vc.Model.Get(id)
}

func (vc *VouchersController) Store(c *gin.Context) {
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
	// Check if ID is already known on data
	_, err = vc.Model.Get(voucherData.ID)
	if err == nil {
		return nil, config.ErrorAlreadyExists
	}
	voucherData.Status = "PENDING"
	err = vc.Model.Update(voucherData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	// Store ID on user information
	err = vc.UserModel.AddVoucher(voucherData.UID, voucherData.ID)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
