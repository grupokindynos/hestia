package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
)

/*

	CardsController is a safe-access query for cards on Firestore Database
	Database Structure:

	cards/
		UID/
          	cards -> Array of CardCodes

	cardsIndex/
		CardCode/
			carddata

	pins/
		CardCode/
			pinhash

*/

type CardsController struct {
	Model *models.CardsModel
}

// User methods

func (cc *CardsController) GetUserAll(uid string, params gin.Params) (interface{}, error) {
	return nil, nil
}

func (cc *CardsController) GetUserSingle(uid string, params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}

// Admin methods

func (cc *CardsController) GetAll(params gin.Params) (interface{}, error) {
	return nil, nil
}

func (cc *CardsController) GetSingle(params gin.Params) (interface{}, error) {
	return models.Shift{}, nil
}

func (cc *CardsController) Store(params gin.Params) (interface{}, error) {
	return "", nil
}
