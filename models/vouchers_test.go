package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVouchersModel_Update(t *testing.T) {
	err := vouchersModel.Update(TestVoucher)
	assert.Nil(t, err)
}

func TestVouchersModel_Get(t *testing.T) {
	newVoucher, err := vouchersModel.Get(TestVoucher.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestVoucher, newVoucher)
}

func TestVouchersModel_GetAll(t *testing.T) {
	vouchers, err := vouchersModel.GetAll("all", "")
	assert.Nil(t, err)
	assert.NotZero(t, len(vouchers))
	assert.IsType(t, []hestia.Voucher{}, vouchers)
}
