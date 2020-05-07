package controllers

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
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

func (gc *GlobalConfigController) GetConfig(userData hestia.User, params Params) (interface{}, error) {
	configData, err := gc.Model.GetConfigData()
	if err != nil {
		return nil, errors.ErrorCoinDataGet
	}
	return configData, nil
}

func (gc *GlobalConfigController) GetConfigMicroservice(c *gin.Context) {
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	configData, err := gc.Model.GetConfigData()
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), configData, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (gc *GlobalConfigController) UpdateConfigData(userData hestia.User, params Params) (interface{}, error) {
	var newConfig hestia.Config
	err := json.Unmarshal(params.Body, &newConfig)
	if err != nil {
		return nil, errors.ErrorUnmarshal
	}
	err = gc.Model.UpdateConfigData(newConfig)
	if err != nil {
		return nil, errors.ErrorDBStore
	}
	return true, nil
}
