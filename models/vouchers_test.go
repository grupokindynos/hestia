package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
