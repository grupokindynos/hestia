package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	VouchersController is a safe-access query for vouchers on Firestore Database
	Database Structure:

	vouchers/
		UID/
          	vouchers -> Array of VouchersIDs

	voucherIndex/
		VoucherID/
			voucherData

*/

type VouchersController struct {
	Model *models.VouchersModel
}

// User methods

func (vc *VouchersController) GetUserAll(uid string, params gin.Params) (interface{}, error) {
	return nil, nil
}

func (vc *VouchersController) GetUserSingle(uid string, params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}

func (vc *VouchersController) Store(uid string, params gin.Params) (interface{}, error) {
	return "", nil
}

// Admin methods

func (vc *VouchersController) GetAll(params gin.Params) (interface{}, error) {
	return nil, nil
}

func (vc *VouchersController) GetSingle(params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}
