package controllers

import (
	"github.com/grupokindynos/hestia/services"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

type Shift struct {
	PaymentAddress string    `bson:"payment_address" json:"payment_address"`
	PaymentAmount  string    `bson:"payment_amount" json:"payment_amount"`
	PaymentCoin    string    `bson:"payment_coin" json:"payment_coin"`
	PaymentTxid    string    `bson:"payment_txid" json:"payment_txid"`
	Rate           float64   `bson:"rate" json:"rate"`
	ID             string    `bson:"id" json:"id"`
	Status         string    `bson:"status" json:"status"`
	Time           time.Time `bson:"time" json:"time"`
	ToAddress      string    `bson:"to_address" json:"to_address"`
	ToAmount       string    `bson:"to_amount" json:"to_amount"`
	ToCoin         string    `bson:"to_coin" json:"to_coin"`
}

type ShiftsController struct {
	DB   *mongo.Database
	Obol *services.ObolService
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
