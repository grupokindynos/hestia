package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
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
	userInfo, err := dc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []models.Deposit
	for _, id := range userInfo.Deposits {
		obj, err := dc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (dc *DepositsController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("depositid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	userInfo, err := dc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Deposits, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return dc.Model.Get(id)
}

func (dc *DepositsController) Store(userData models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (dc *DepositsController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	objs, err := dc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorAllError
	}
	return objs, nil
}

func (dc *DepositsController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("depositid")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return dc.Model.Get(id)
}

func (dc *DepositsController) Update(userData models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}
