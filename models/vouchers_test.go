package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestVoucher = Voucher{
	ID:         "TEST-VOUCHER",
	UID:        "XYZ12345678910",
	VoucherID:  "FAKE-VOUCHER",
	VariantID:  "FAKE-VARIANT",
	FiatAmount: "100",
	Name:       "TEST-VOUCHER",
	PaymentData: Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	BitcouPaymentData: Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	RedeemCode: "FAKE-REDEEM",
	Status:     "COMPLETED",
	Timestamp:  "00000000000",
}

func TestVouchersModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := VouchersModel{
		Db:         db,
		Collection: "vouchers",
	}
	err = model.Update(TestVoucher)
	assert.Nil(t, err)
}

func TestVouchersModel_Get(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := VouchersModel{
		Db:         db,
		Collection: "vouchers",
	}
	newVoucher, err := model.Get(TestVoucher.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestVoucher, newVoucher)
}

func TestVouchersModel_GetAll(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := VouchersModel{
		Db:         db,
		Collection: "vouchers",
	}
	vouchers, err := model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(vouchers))
	assert.IsType(t, []Voucher{}, vouchers)
}
