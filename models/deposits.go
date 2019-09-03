package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Deposit struct{}

type DepositsModel struct {
	Db *mongo.Database
}

func (m *DepositsModel) GetDeposit(depositid string) (deposit Deposit, err error) {
	return deposit, err
}

func (m *DepositsModel) StoreDeposit(uid string, deposit Deposit) error {
	return nil
}
