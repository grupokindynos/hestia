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
	collection := m.Db.Collection("shifts")
	filter := bson.M{"_id": shiftid}
	err = collection.FindOne(ctx, filter).Decode(&shift)
	return shift, err
}

func (m *ShiftModel) UpdateShift(uid string, shift Shift) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftIndexColl := m.Db.Collection("shifts")
	shiftIDFilter := bson.M{"_id": shift.ID}
	upsert := true
	_, err := shiftIndexColl.UpdateOne(ctx, shiftIDFilter, bson.D{{Key: "$set", Value: shift}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *ShiftModel) GetAll() (shifts []Shift, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.Db.Collection("shifts")
	curr, err := collection.Find(ctx, nil)
	if err != nil {
		return shifts, err
	}
	for curr.Next(ctx) {
		var shift Shift
		err := curr.Decode(&shift)
		if err != nil {
			return shifts, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, err
}
