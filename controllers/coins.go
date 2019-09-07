package controllers

import (
	"github.com/gin-gonic/gin"
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

// Admin methods

func (cc *CoinsController) GetCoinsAvailability(c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (cc *CoinsController) UpdateCoinsAvailability(c *gin.Context) (interface{}, error) {
	return nil, nil
}
