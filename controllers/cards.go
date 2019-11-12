package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/utils"
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

func (cc *CardsController) GetAll(userData hestia.User, params Params) (interface{}, error) {
	if params.Admin {
		return cc.Model.GetAll()
	}
	userInfo, err := cc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var Array []hestia.Card
	for _, id := range userInfo.Cards {
		obj, err := cc.Model.Get(id)
		if err != nil {
			return nil, errors.ErrorNotFound
		}
		Array = append(Array, obj)
	}
	return Array, nil
}

func (cc *CardsController) GetSingle(userData hestia.User, params Params) (interface{}, error) {
	if params.CardCode == "" {
		return nil, errors.ErrorMissingID
	}
	if params.Admin {
		return cc.Model.Get(params.CardCode)
	}
	userInfo, err := cc.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Cards, params.CardCode) {
		return nil, errors.ErrorInfoDontMatchUser
	}
	return cc.Model.Get(params.CardCode)
}

func (cc *CardsController) Store(c *gin.Context) {
	// Catch the request jwe
	var ReqBody hestia.BodyReq
	err := c.BindJSON(&ReqBody)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	// Verify Signature
	// TODO here we need to use Cards Microservice signature
	rawBytes, err := jwt.DecodeJWS(ReqBody.Payload, os.Getenv(""))
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDecryptJWE, c)
		return
	}
	// Try to unmarshal the information of the payload
	var cardData hestia.Card
	err = json.Unmarshal(rawBytes, &cardData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	_, err = cc.Model.Get(cardData.CardCode)
	if err == nil {
		responses.GlobalResponseError(nil, errors.ErrorAlreadyExists, c)
		return
	}
	err = cc.Model.Update(cardData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	// Store ID on user information
	err = cc.UserModel.AddVoucher(cardData.UID, cardData.CardCode)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	response, err := jwt.EncodeJWS(cardData.CardCode, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseError(response, err, c)
	return
}
