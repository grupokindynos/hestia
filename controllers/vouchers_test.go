package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestVouchersController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	vouchers, err := vouchersCtrl.GetAll(models.TestUser, c, false, "all")
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
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "voucherid", Value: models.TestVoucher.ID}}
	voucher, err := vouchersCtrl.GetSingle(models.TestUser, c, false, "all")
	assert.Nil(t, err)
	assert.IsType(t, hestia.Voucher{}, voucher)
	assert.Equal(t, models.TestVoucher, voucher)
}

func TestVouchersController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	vouchers, err := vouchersCtrl.GetAll(models.TestUser, c, true, "all")
	assert.Nil(t, err)
	var voucherArray []hestia.Voucher
	voucherBytes, err := json.Marshal(vouchers)
	assert.Nil(t, err)
	err = json.Unmarshal(voucherBytes, &voucherArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Voucher{}, vouchers)
	assert.Equal(t, models.TestVoucher, voucherArray[0])
}

func TestVouchersController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "voucherid", Value: models.TestVoucher.ID}}
	voucher, err := vouchersCtrl.GetSingle(models.TestUser, c, true, "all")
	assert.Nil(t, err)
	assert.IsType(t, hestia.Voucher{}, voucher)
	assert.Equal(t, models.TestVoucher, voucher)
}
