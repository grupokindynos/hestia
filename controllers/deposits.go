package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	DepositsController is a safe-access query for deposits on Firestore Database
	Database Structure:

	deposits/
		DepositID/
			depositData

*/

type DepositsController struct {
	Model     *models.DepositsModel
	UserModel *models.UsersModel
}

// User methods

func (dc *DepositsController) GetUserAll(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (dc *DepositsController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}

func (dc *DepositsController) Store(userData models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (dc *DepositsController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (dc *DepositsController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}
