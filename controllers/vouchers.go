package controllers

import (
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
	return "", nil
}

// Admin methods

func (vc *VouchersController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (vc *VouchersController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("voucherid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return vc.Model.Get(id)
}
