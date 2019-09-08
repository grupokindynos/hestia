package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Voucher struct {
	ID                string  `bson:"id" json:"id"`
	UID               string  `bson:"uid" json:"uid"`
	VoucherID         string  `bson:"voucher_id" json:"voucher_id"`
	VariantID         string  `bson:"variant_id" json:"variant_id"`
	FiatAmount        string  `bson:"fiat_amount" json:"fiat_amount"`
	Name              string  `bson:"name" json:"name"`
	PaymentData       Payment `bson:"payment_data" json:"payment_data"`
	BitcouPaymentData Payment `bson:"bitcou_payment_data" json:"bitcou_payment_data"`
	RedeemCode        string  `bson:"redeem_code" json:"redeem_code"`
	Status            string  `bson:"status" json:"status"`
	Timestamp         string  `bson:"timestamp" json:"timestamp"`
	TxnID             string  `bson:"txn_id" json:"txn_id"`
}

type VouchersModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *VouchersModel) Get(id string) (voucher Voucher, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&voucher)
	return voucher, err
}

func (m *VouchersModel) Update(id string, voucher Voucher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: voucher}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *VouchersModel) GetAll() (vouchers []Voucher, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, err := col.Find(ctx, nil)
	if err != nil {
		return vouchers, err
	}
	for curr.Next(ctx) {
		var voucher Voucher
		err := curr.Decode(&voucher)
		if err != nil {
			return vouchers, err
		}
		vouchers = append(vouchers, voucher)
	}
	return vouchers, err
}
