package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/jws"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
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

func (cc *CardsController) GetAll(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return cc.Model.GetAll()
	}
	userInfo, err := cc.UserModel.Get(userData.ID)
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

func (cc *CardsController) GetSingle(userData models.User, c *gin.Context, admin bool) (interface{}, error) {
	id, ok := c.Params.Get("cardcode")
	if !ok {
		return nil, config.ErrorMissingID
	}
	if admin {
		return cc.Model.Get(id)
	}
	userInfo, err := cc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Cards, id) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return cc.Model.Get(id)
}

func (cc *CardsController) Store(c *gin.Context) {
	return
}
