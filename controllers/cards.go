package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
)

/*

	CardsController is a safe-access query for cards on Firestore Database
	Database Structure:

	cards/
		cardcode/
          	carddata

	pins/
		CardCode/
			pinhash

*/

type CardsController struct {
	Model     *models.CardsModel
	UserModel *models.UsersModel
}

// User methods

func (cc *CardsController) GetUserAll(userData models.User, c *gin.Context) (interface{}, error) {
	userInfo, err := cc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []models.Card
	for _, id := range userInfo.Cards {
		obj, err := cc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (cc *CardsController) GetUserSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("cardcode")
	if !ok {
		return nil, config.ErrorMissingID
	}
	userInfo, err := cc.UserModel.GetUserInformation(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Cards, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return cc.Model.Get(id)
}

// Admin methods

func (cc *CardsController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	objs, err := cc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorAllError
	}
	return objs, nil
}

func (cc *CardsController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	id, ok := c.Params.Get("cardcode")
	if !ok {
		return nil, config.ErrorMissingID
	}
	return cc.Model.Get(id)
}

func (cc *CardsController) Store(userData models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}
