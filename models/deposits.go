package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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
	Db         *mongo.Database
	Collection string
}

func (m *DepositsModel) Get(id string) (deposit Deposit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&deposit)
	return deposit, err
}

func (m *DepositsModel) Update(deposit Deposit) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": deposit.ID}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: deposit}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *DepositsModel) GetAll() (deposits []Deposit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, err := col.Find(ctx, bson.M{})
	if err != nil {
		return deposits, err
	}
	for curr.Next(ctx) {
		var deposit Deposit
		err := curr.Decode(&deposit)
		if err != nil {
			return deposits, err
		}
		deposits = append(deposits, deposit)
	}
	return deposits, err
}
