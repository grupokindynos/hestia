package models

import (
	"context"
	"github.com/grupokindynos/common/hestia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CoinsModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *CoinsModel) GetCoinsData() ([]hestia.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	var CoinData []hestia.Coin
	cursor, _ := col.Find(ctx, bson.M{})
	for cursor.Next(ctx) {
		var coinProp hestia.Coin
		_ = cursor.Decode(&coinProp)
		CoinData = append(CoinData, coinProp)
	}
	return CoinData, nil
}

func (m *CoinsModel) UpdateCoinsData(Coins []hestia.Coin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	for _, coin := range Coins {
		filter := bson.M{"_id": coin.Ticker}
		upsert := true
		_, _ = col.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: coin}}, &options.UpdateOptions{Upsert: &upsert})
	}
	return nil
}
