package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/utils"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"os"
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

func (cc *CardsController) GetAll(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
	if admin {
		return cc.Model.GetAll()
	}
	userInfo, err := cc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var Array []hestia.Card
	for _, id := range userInfo.Cards {
		obj, err := cc.Model.Get(id)
		if err != nil {
			return nil, config.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (cc *CardsController) GetSingle(userData hestia.User, c *gin.Context, admin bool) (interface{}, error) {
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
	// Catch the request jwe
	var ReqBody hestia.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Cards Microservice signature
	rawBytes, err := jwt.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var cardData hestia.Card
	err = json.Unmarshal(rawBytes, &cardData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorUnmarshal, c)
		return
	}
	// Hash the PaymentTxID as the ID
	cardData.CardCode = fmt.Sprintf("%x", sha256.Sum256([]byte(cardData.CardNumber)))
	// Check if ID is already known on data
	_, err = cc.Model.Get(cardData.CardCode)
	if err == nil {
		responses.GlobalResponseError(nil, config.ErrorAlreadyExists, c)
		return
	}
	err = cc.Model.Update(cardData)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = cc.UserModel.AddVoucher(cardData.UID, cardData.CardCode)
	if err != nil {
		responses.GlobalResponseError(nil, config.ErrorDBStore, c)
		return
	}
	response, err := jwt.EncodeJWS(cardData.CardCode, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseError(response, err, c)
	return
}
