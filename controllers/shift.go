package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services"
	"github.com/grupokindynos/hestia/utils"
)

/*

	ShiftController is a safe-access query for shifts on Firestore Database
	Database Structure:

	shift/
		ShiftID/
			shiftData

*/

type ShiftsController struct {
	Obol      *services.ObolService
	Model     *models.ShiftModel
	UserModel *models.UsersModel
}

// User methods

func (sc *ShiftsController) GetUserAll(userInfo models.User, c *gin.Context) (interface{}, error) {
	userInfo, err := sc.UserModel.GetUserInformation(userInfo.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	var ShiftsArray []models.Shift
	for _, shiftID := range userInfo.Shifts {
		shift, err := sc.Model.GetShift(shiftID)
		if err != nil {
			return nil, config.ErrorShiftNotFound
		}
		ShiftsArray = append(ShiftsArray, shift)
	}
	return ShiftsArray, nil
}

func (sc *ShiftsController) GetUserSingle(userInfo models.User, c *gin.Context) (interface{}, error) {
	shiftID, ok := c.Params.Get("shiftid")
	if !ok {
		return nil, config.ErrorMissingShiftID
	}
	userInfo, err := sc.UserModel.GetUserInformation(userInfo.ID)
	if err != nil {
		return nil, config.ErrorNoUserInformation
	}
	if !utils.Contains(userInfo.Shifts, shiftID) {
		return nil, config.ErrorInfoDontMatchUser
	}
	return sc.Model.GetShift(shiftID)
}

func (sc *ShiftsController) Store(userInfo models.User, c *gin.Context) (interface{}, error) {
	return "", nil
}

// Admin methods

func (sc *ShiftsController) GetAll(userData models.User, c *gin.Context) (interface{}, error) {
	shifts, err := sc.Model.GetAll()
	if err != nil {
		return nil, config.ErrorShiftsAllError
	}
	return shifts, nil
}

func (sc *ShiftsController) GetSingle(userData models.User, c *gin.Context) (interface{}, error) {
	shiftID, ok := c.Params.Get("shiftid")
	if !ok {
		return nil, config.ErrorMissingShiftID
	}
	return sc.Model.GetShift(shiftID)
}

func (sc *ShiftsController) Update(userData models.User, c *gin.Context) (interface{}, error) {
	return models.Shift{}, nil
}
