package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/config"
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

func (cc *CoinsController) GetCoinsAvailability(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		return nil, config.ErrorCoinDataGet
	}
	return coins, nil
}

func (cc *CoinsController) GetCoinsAvailabilityMicroService(c *gin.Context) {
	headerSignature := c.GetHeader("service")
	if headerSignature == "" {
		responses.GlobalResponseNoAuth(c)
		return
	}
	decodedHeader, err := jwt.DecodeJWSNoVerify(headerSignature)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	var serviceStr string
	err = json.Unmarshal(decodedHeader, &serviceStr)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Check which service the request is announcing
	var pubKey string
	switch serviceStr {
	case "ladon":
		pubKey = os.Getenv("LADON_PUBLIC_KEY")
	case "tyche":
		pubKey = os.Getenv("TYCHE_PUBLIC_KEY")
	case "adrestia":
		pubKey = os.Getenv("ADRESTIA_PUBLIC_KEY")
	default:
		responses.GlobalResponseNoAuth(c)
		return
	}
	valid, _ := mvt.VerifyMVTToken(headerSignature, nil, pubKey, os.Getenv("MASTER_PASSWORD"))
	if !valid {
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
