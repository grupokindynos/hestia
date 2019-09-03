package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Voucher struct{}

type VouchersModel struct {
	Db *mongo.Database
}

func (m *VouchersModel) GetVoucher(voucherid string) (voucher Voucher, err error) {
	return voucher, err
}

func (m *VouchersModel) StoreVoucher(uid string, voucher Voucher) error {
	return nil
}
