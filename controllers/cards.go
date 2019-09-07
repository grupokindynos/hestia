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

func (cc *CardsController) GetUserAll(userInfo models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (cc *CardsController) GetUserSingle(userInfo models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}

// Admin methods

func (cc *CardsController) GetAll(c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (cc *CardsController) GetSingle(c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}

func (cc *CardsController) Store(c *gin.Context) (interface{}, error) {
	return "", nil
}
