package models

import (
	"context"
	"github.com/grupokindynos/common/hestia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DepositsModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *DepositsModel) Get(id string) (deposit hestia.Deposit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&deposit)
	return deposit, err
}

func (m *DepositsModel) Update(deposit hestia.Deposit) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": deposit.ID}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: deposit}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *DepositsModel) GetAll() (deposits []hestia.Deposit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, _ := col.Find(ctx, bson.M{})
	for curr.Next(ctx) {
		var deposit hestia.Deposit
		_ = curr.Decode(&deposit)
		deposits = append(deposits, deposit)
	}
	return deposits, nil
}
