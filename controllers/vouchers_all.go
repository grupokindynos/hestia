package controllers

import (
	"fmt"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/shopspring/decimal"
	"strconv"
)

type VouchersAllController struct {
	UserModel       *models.UsersModel
	VouchersModel   *models.VouchersModel
	VouchersV2Model *models.VouchersModelV2
}

func voucherToLightVoucher(voucher hestia.Voucher) hestia.LightVoucher {
	dec := decimal.NewFromInt(voucher.PaymentData.Amount).Div(decimal.NewFromInt(1e8))
	amount, _ := dec.Float64()
	variantId, _ := strconv.Atoi(voucher.VariantID)
	return hestia.LightVoucher{
		Id:             voucher.ID,
		VoucherId:      voucher.VoucherID,
		Name:           voucher.Name,
		Timestamp:      voucher.Timestamp,
		Amount:         amount,
		PaymentTxId:    voucher.PaymentData.Txid,
		PaymentCoin:    voucher.PaymentData.Coin,
		RefundTxId:     "", // VouchersV1 don't have refund info.
		Status:         voucher.Status,
		ProviderId:     fmt.Sprintf("%d", voucher.ProviderId),
		ShippingMethod: hestia.VoucherShippingMethodApi, // vouchersv1 don't have other option.
		VariantId:      variantId,
		RedeemCode:     voucher.RedeemCode,
	}
}

func voucherV2toLightVoucher(voucher hestia.VoucherV2) hestia.LightVoucher {
	dec := decimal.NewFromInt(voucher.UserPayment.Amount).Div(decimal.NewFromInt(1e8))
	amount, _ := dec.Float64()
	return hestia.LightVoucher{
		Id:             voucher.Id,
		VoucherId:      voucher.VoucherId,
		Name:           voucher.VoucherName,
		Timestamp:      voucher.CreatedTime,
		Amount:         amount,
		PaymentTxId:    voucher.UserPayment.Txid,
		PaymentCoin:    voucher.UserPayment.Coin,
		RefundTxId:     voucher.RefundTxId,
		Status:         hestia.GetVoucherStatusV2String(voucher.Status),
		ProviderId:     voucher.ProviderId,
		ShippingMethod: voucher.ShippingMethod,
		VariantId:      voucher.VariantId,
		RedeemCode:     voucher.RedeemCode,
	}
}

func (va *VouchersAllController) GetVouchersHistory(userData hestia.User, _ Params) (interface{}, error) {
	userInfo, err := va.UserModel.Get(userData.ID)
	if err != nil {
		return nil, errors.ErrorNoUserInformation
	}
	var lightVouchers []hestia.LightVoucher

	for _, id := range userInfo.Vouchers {
		obj, err := va.VouchersModel.Get(id)
		if err != nil {
			continue
		}
		voucher := voucherToLightVoucher(obj)
		lightVouchers = append(lightVouchers, voucher)
	}

	for _, id := range userInfo.VouchersV2 {
		obj, err := va.VouchersV2Model.Get(id)
		if err != nil {
			continue
		}
		voucher := voucherV2toLightVoucher(obj)
		lightVouchers = append(lightVouchers, voucher)
	}
	return lightVouchers, nil
}
