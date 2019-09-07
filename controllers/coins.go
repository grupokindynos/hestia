package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/utils"
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

// Admin methods

func (cc *CoinsController) GetCoinsAvailability(userData models.User, c *gin.Context) (interface{}, error) {
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		return nil, config.ErrorCoinDataGet
	}
	return coins, nil
}

func (cc *CoinsController) UpdateCoinsAvailability(userData models.User, c *gin.Context) (interface{}, error) {
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	rawBytes, err := utils.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, config.ErrorDecryptJWE
	}
	var newCoinsData []models.Coin
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
