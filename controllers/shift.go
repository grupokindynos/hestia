package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services"
)

/*

	ShiftController is a safe-access query for shifts on Firestore Database
	Database Structure:

	shift/
		UID/
          	shifts -> Array of ShiftsID

	shiftIndex/
		ShiftID/
			shiftData

*/

type ShiftsController struct {
	Obol  *services.ObolService
	Model *models.ShiftModel
}

// User methods

func (sc *ShiftsController) GetUserAll(userInfo models.User, c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (sc *ShiftsController) GetUserSingle(userInfo models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}

func (sc *ShiftsController) Store(userInfo models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (sc *ShiftsController) GetAll(c *gin.Context) (interface{}, error) {
	return nil, nil
}

func (sc *ShiftsController) GetSingle(c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}
