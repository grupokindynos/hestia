package models

import (
	"go.mongodb.org/mongo-driver/mongo"
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
	Db *mongo.Database
}
