package controllers

import "cloud.google.com/go/firestore"

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

type Shift struct{}

type ShiftsController struct {
	DB *firestore.Client
}

func (sc *ShiftsController) GetShift(shiftid string) (shift Shift, err error) {
	return Shift{}, nil
}

func (sc *ShiftsController) GetUserShifts(uid string) (shifts []Shift, err error) {
	return nil, nil
}

func (sc *ShiftsController) GetUserShiftsIDs(uid string) (shiftsIDs []string, err error) {
	return nil, nil
}
