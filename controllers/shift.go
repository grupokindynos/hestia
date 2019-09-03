package controllers

import (
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services"
	"go.mongodb.org/mongo-driver/mongo"
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
	DB   *mongo.Database
	Obol *services.ObolService
}

func (sc *ShiftsController) GetShift(shiftid string) (shift models.Shift, err error) {
	return models.Shift{}, nil
}

func (sc *ShiftsController) GetUserShifts(uid string) (shifts []models.Shift, err error) {
	return nil, nil
}

func (sc *ShiftsController) GetUserShiftsIDs(uid string) (shiftsIDs []string, err error) {
	return nil, nil
}
