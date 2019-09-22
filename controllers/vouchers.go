package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jws"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"os"
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

func (vc *VouchersController) GetAll(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return vc.Model.GetAll()
	}
	userInfo, err := vc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []hestia.Voucher
	for _, id := range userInfo.Vouchers {
		obj, err := vc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (vc *VouchersController) GetSingle(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
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
	if !utils.Contains(userInfo.Vouchers, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return vc.Model.Get(id)
}

func (vc *VouchersController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Vouchers Microservice signature
	rawBytes, err := jws.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var voucherData hestia.Voucher
	err = json.Unmarshal(rawBytes, &voucherData)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	voucherData.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(voucherData.PaymentData.Txid)))
	// Check if ID is already known on data
	_, err = vc.Model.Get(voucherData.ID)
	if err == nil {
		config.GlobalResponseError(nil, config.ErrorAlreadyExists, c)
		return
	}
	voucherData.Status = "PENDING"
	err = vc.Model.Update(voucherData)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = vc.UserModel.AddVoucher(voucherData.UID, voucherData.ID)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	response, err := jws.EncodeJWS(voucherData.ID, os.Getenv("HESTIA_PRIVATE_KEY"))
	config.GlobalResponseError(response, err, c)
	return
}
