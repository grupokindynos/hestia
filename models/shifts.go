package models

import (
	"context"
	"github.com/grupokindynos/common/hestia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ShiftModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *ShiftModel) Get(id string) (shift hestia.Shift, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&shift)
	return shift, err
}

func (m *ShiftModel) Update(shift hestia.Shift) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": shift.ID}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: shift}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *ShiftModel) GetAll() (shifts []hestia.Shift, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, _ := col.Find(ctx, bson.M{})
	for curr.Next(ctx) {
		var shift hestia.Shift
		_ = curr.Decode(&shift)
		shifts = append(shifts, shift)
	}
	return shifts, nil
}
