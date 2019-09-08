package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Coin struct {
	Ticker            string   `bson:"ticker" json:"ticker"`
	ShiftAvailable    bool     `bson:"shift_available" json:"shift_available"`
	DepositAvailable  bool     `bson:"deposit_available" json:"deposit_available"`
	VouchersAvailable bool     `bson:"vouchers_available" json:"vouchers_available"`
	OrdersAvailable   bool     `bson:"orders_available" json:"orders_available"`
	Balances          Balances `bson:"balances" json:"balances"`
}

type Balances struct {
	HotWallet float64 `bson:"hot_wallet" json:"hot_wallet"`
	Exchanges float64 `bson:"exchanges" json:"exchanges"`
}

type CoinsModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *CoinsModel) GetCoinsData() ([]Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	var CoinData []Coin
	cursor, err := col.Find(ctx, nil)
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

func (m *CoinsModel) UpdateCoinsData(Coins []Coin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	for _, coin := range Coins {
		filter := bson.M{"_id": coin.Ticker}
		upsert := true
		_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: coin}}, &options.UpdateOptions{Upsert: &upsert})
		if err != nil {
			return err
		}
	}
	return nil
}
