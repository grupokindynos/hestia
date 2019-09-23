package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
)

/*

	CoinsController is a safe-access query for cards on Firestore Database
	Database Structure:

	coins/
		TICKER/
          	coinAvailability

*/

type CoinsController struct {
	Model *models.CoinsModel
}

func (cc *CoinsController) GetCoinsAvailability(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		return nil, config.ErrorCoinDataGet
	}
	return coins, nil
}

func (cc *CoinsController) UpdateCoinsAvailability(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	rawBytes, err := jwt.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, config.ErrorDecryptJWE
	}
	var newCoinsData []hestia.Coin
	err = json.Unmarshal(rawBytes, &newCoinsData)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	err = cc.Model.UpdateCoinsData(newCoinsData)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
