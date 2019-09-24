package models

import (
	"context"
	"github.com/grupokindynos/common/hestia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type VouchersModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *VouchersModel) Get(id string) (voucher hestia.Voucher, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&voucher)
	return voucher, err
}

func (m *VouchersModel) Update(voucher hestia.Voucher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": voucher.ID}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: voucher}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *VouchersModel) GetAll(filter string) (vouchers []hestia.Voucher, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, _ := col.Find(ctx, bson.M{})
	for curr.Next(ctx) {
		var voucher hestia.Voucher
		_ = curr.Decode(&voucher)
		if filter != "all" {
			if voucher.Status == strings.ToUpper(filter) {
				vouchers = append(vouchers, voucher)
			}
		} else {
			vouchers = append(vouchers, voucher)
		}
	}
	return vouchers, nil
}
