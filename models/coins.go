package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Coin struct {
	ShiftAvailable    bool `bson:"shift_available" json:"shift_available"`
	DepositAvailable  bool `bson:"deposit_available" json:"deposit_available"`
	VouchersAvailable bool `bson:"vouchers_available" json:"vouchers_available"`
	OrdersAvailable   bool `bson:"orders_available" json:"orders_available"`
}

type CoinsModel struct {
	Db *mongo.Database
}
