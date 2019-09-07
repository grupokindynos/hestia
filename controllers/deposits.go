package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	DepositsController is a safe-access query for deposits on Firestore Database
	Database Structure:

	deposits/
		UID/
          	deposit -> Array of DepositIDs

	depositIndex/
		DepositID/
			depositData

*/

type DepositsController struct {
	Model *models.DepositsModel
}

// User methods

func (dc *DepositsController) GetUserAll(uid string, params gin.Params) (interface{}, error) {
	return nil, nil
}

func (dc *DepositsController) GetUserSingle(uid string, params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}

func (dc *DepositsController) Store(uid string, params gin.Params) (interface{}, error) {
	return "", nil
}

// Admin methods

func (dc *DepositsController) GetAll(params gin.Params) (interface{}, error) {
	return nil, nil
}

func (dc *DepositsController) GetSingle(params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}
