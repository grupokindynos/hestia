package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Shift struct {
	ID         string  `bson:"id" json:"id"`
	Status     string  `bson:"status" json:"status"`
	Timestamp  string  `bson:"timestamp" json:"timestamp"`
	Payment    Payment `bson:"payment" json:"payment"`
	Conversion Payment `bson:"conversion" json:"conversion"`
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
