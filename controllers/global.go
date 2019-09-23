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

	GlobalConfigController is a safe-access query for cards on Firestore Database
	Database Structure:

	config/
		shifts/props
		deposits/props
		vouchers/props
		orders/props

*/

type GlobalConfigController struct {
	Model *models.GlobalConfigModel
}

func (gc *GlobalConfigController) GetConfig(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	coins, err := gc.Model.GetConfigData()
	if err != nil {
		return nil, config.ErrorCoinDataGet
	}
	return coins, nil
}

func (gc *GlobalConfigController) UpdateConfigData(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	var ReqBody models.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	rawBytes, err := jwt.DecryptJWE(userData.ID, ReqBody.Payload)
	if err != nil {
		return nil, config.ErrorDecryptJWE
	}
	var newConfig hestia.Config
	err = json.Unmarshal(rawBytes, &newConfig)
	if err != nil {
		return nil, config.ErrorUnmarshal
	}
	err = gc.Model.UpdateConfigData(newConfig)
	if err != nil {
		return nil, config.ErrorDBStore
	}
	return true, nil
}
