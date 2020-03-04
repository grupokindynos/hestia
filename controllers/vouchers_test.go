package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVouchersController_GetUserAll(t *testing.T) {
	vouchers, err := vouchersCtrl.GetAll(models.TestUser, TestParams)
	assert.Nil(t, err)
	var vouchersArray []hestia.Voucher
	voucherBytes, err := json.Marshal(vouchers)
	assert.Nil(t, err)
	err = json.Unmarshal(voucherBytes, &vouchersArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Voucher{}, vouchers)
	assert.Equal(t, models.TestVoucher, vouchersArray[0])
}

func TestVouchersController_GetUserSingle(t *testing.T) {
	voucher, err := vouchersCtrl.GetSingle(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Voucher{}, voucher)
	assert.Equal(t, models.TestVoucher, voucher)
}

func TestVouchersController_GetAll(t *testing.T) {
	vouchers, err := vouchersCtrl.GetAll(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	var voucherArray []hestia.Voucher
	voucherBytes, err := json.Marshal(vouchers)
	assert.Nil(t, err)
	err = json.Unmarshal(voucherBytes, &voucherArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Voucher{}, vouchers)
}

func TestVouchersController_GetSingle(t *testing.T) {
	voucher, err := vouchersCtrl.GetSingle(models.TestUser, TestParamsAdmin)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Voucher{}, voucher)
	assert.Equal(t, models.TestVoucher, voucher)
}
