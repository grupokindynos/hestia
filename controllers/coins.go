package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
	"os"
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

func (cc *CoinsController) GetCoinsAvailability(userData hestia.User, params Params) (interface{}, error) {
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		return nil, errors.ErrorCoinDataGet
	}
	return coins, nil
}

func (cc *CoinsController) GetCoinsAvailabilityMicroService(c *gin.Context) {
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), coins, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (cc *CoinsController) UpdateCoinsAvailability(userData hestia.User, params Params) (interface{}, error) {
	var newCoinsData []hestia.Coin
	err := json.Unmarshal(params.Body, &newCoinsData)
	if err != nil {
		return nil, errors.ErrorUnmarshal
	}
	err = cc.Model.UpdateCoinsData(newCoinsData)
	if err != nil {
		return nil, errors.ErrorDBStore
	}
	return true, nil
}
