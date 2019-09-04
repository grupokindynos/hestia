package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Deposit struct {
	ID           string  `bson:"id" json:"id"`
	UID          string  `bson:"uid" json:"uid"`
	Payment      Payment `bson:"payment" json:"payment"`
	AmountInPeso string  `bson:"amount_in_peso" json:"amount_in_peso"`
	CardCode     string  `bson:"card_code" json:"card_code"`
	Status       string  `bson:"status" json:"status"`
	Timestamp    string  `bson:"timestamp" json:"timestamp"`
}

type DepositsModel struct {
	Db *mongo.Database
}

func (m *DepositsModel) GetDeposit(depositid string) (deposit Deposit, err error) {
	return deposit, err
}

func (m *DepositsModel) StoreDeposit(uid string, deposit Deposit) error {
	return nil
}
