package test

import (
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVouchersModel_Update(t *testing.T) {
	err := vouchersCtrl.Model.Update(TestVoucher)
	assert.Nil(t, err)
}

func TestVouchersModel_Get(t *testing.T) {
	newVoucher, err := vouchersCtrl.Model.Get(TestVoucher.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestVoucher, newVoucher)
}

func TestVouchersModel_GetAll(t *testing.T) {
	vouchers, err := vouchersCtrl.Model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(vouchers))
	assert.IsType(t, []models.Voucher{}, vouchers)
}
