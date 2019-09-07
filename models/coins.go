package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

func (m *CoinsModel) GetCoinsData() ([]Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.Db.Collection("shiftIndex")
	var CoinData []Coin
	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var coinProp Coin
		err := cursor.Decode(&coinProp)
		if err != nil {
			return nil, err
		}
		CoinData = append(CoinData, coinProp)
	}
	return CoinData, nil
}
