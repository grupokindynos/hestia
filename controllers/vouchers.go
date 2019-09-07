package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	VouchersController is a safe-access query for vouchers on Firestore Database
	Database Structure:

	vouchers/
		VoucherID/
			voucherData

*/

type VouchersController struct {
	Model *models.VouchersModel
}

// User methods

func (vc *VouchersController) GetUserAll(userInfo models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (vc *VouchersController) GetUserSingle(userInfo models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}

func (vc *VouchersController) Store(userInfo models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (vc *VouchersController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (vc *VouchersController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}
