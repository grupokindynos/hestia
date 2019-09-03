package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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

type ShiftModel struct {
	Db *mongo.Database
}

func (m *ShiftModel) GetShift(shiftid string) (shift Shift, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.Db.Collection("shiftIndex")
	filter := bson.M{"_id": shift.ID}
	err = collection.FindOne(ctx, filter).Decode(&shift)
	return shift, err
}

func (m *ShiftModel) StoreShift(uid string, shift Shift) error {

	// Add ShiftID to /shift/uid/shifts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection("shifts")
	uidFilter := bson.M{"_id": uid}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$push", Value: shift.ID}}, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		return err
	}
	// Update /shiftIndex/shiftID/shift
	shiftIndexColl := m.Db.Collection("shiftIndex")
	shiftIDFilter := bson.M{"_id": shift.ID}
	_, err = shiftIndexColl.UpdateOne(ctx, shiftIDFilter, bson.D{{Key: "$set", Value: shift}}, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		return err
	}
	return err
}
