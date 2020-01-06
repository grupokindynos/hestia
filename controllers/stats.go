package controllers

import (
	"errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/obol"
	"github.com/grupokindynos/hestia/models"
	"strconv"
)

type StatsController struct {
	ShiftModel    *models.ShiftModel
	VouchersModel *models.VouchersModel
	DepositsModel *models.DepositsModel
	OrdersModel   *models.OrdersModel
}

var coinRates = make(map[string][]obol.Rate)

// Big TODO this methods can be merged.

func (sc StatsController) GetShiftStats(userData hestia.User, params Params) (interface{}, error) {
	shifts, err := sc.ShiftModel.GetAll("all")
	if err != nil {
		return nil, err
	}
	response := models.ShiftStatsResponse{
		Pending:    0,
		Confirming: 0,
		Confirmed:  0,
		Error:      0,
		Refund:     0,
		Refunded:   0,
		Complete:   0,
		Total:      0,
	}
	for _, shift := range shifts {
		if shift.Status == hestia.GetShiftStatusString(hestia.ShiftStatusComplete) {
			req := obol.ObolRequest{ObolURL: "https://obol.polispay.com"}
			if shift.FeePayment.Txid != "" {
				feeRates, ok := coinRates[shift.FeePayment.Coin]
				if !ok {
					feeRates, err := req.GetCoinRates(shift.FeePayment.Coin)
					if err != nil {
						return nil, err
					}
					coinRates[shift.FeePayment.Coin] = feeRates
				}
				feeRates, err := req.GetCoinRates(shift.FeePayment.Coin)
				if err != nil {
					return nil, err
				}
				var btcFeeRate float64
				for _, rate := range feeRates {
					if rate.Code == "BTC" {
						btcFeeRate = rate.Rate
					}
				}
				response.Volume += btcFeeRate * (float64(shift.FeePayment.Amount) / 1e8)
			}
			paymentRates, ok := coinRates[shift.Payment.Coin]
			if !ok {
				paymentRates, err := req.GetCoinRates(shift.Payment.Coin)
				if err != nil {
					return nil, err
				}
				coinRates[shift.Payment.Coin] = paymentRates
			}

			var btcPaymentRate float64
			for _, rate := range paymentRates {
				if rate.Code == "BTC" {
					btcPaymentRate = rate.Rate
				}
			}
			response.Volume += btcPaymentRate * (float64(shift.Payment.Amount) / 1e8)
		}
		response.Total += 1
		switch shift.Status {
		case hestia.GetShiftStatusString(hestia.ShiftStatusPending):
			response.Pending += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusConfirming):
			response.Confirming += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusConfirmed):
			response.Confirmed += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusError):
			response.Error += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusRefund):
			response.Refund += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusRefunded):
			response.Refunded += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusComplete):
			response.Complete += 1
		}
	}
	return response, nil
}

func (sc StatsController) GetVoucherStats(userData hestia.User, params Params) (interface{}, error) {
	vouchers, err := sc.VouchersModel.GetAll("all")
	if err != nil {
		return nil, err
	}
	response := models.VoucherStatsResponse{
		Pending:           0,
		Confirming:        0,
		Confirmed:         0,
		AwaitingProvider:  0,
		Error:             0,
		RefundFee:         0,
		RefundTotal:       0,
		RefundedPartially: 0,
		Refunded:          0,
		Complete:          0,
		Total:             0,
	}
	for _, voucher := range vouchers {
		if voucher.Status == hestia.GetVoucherStatusString(hestia.VoucherStatusComplete) {
			response.Volume += float64(voucher.BitcouPaymentData.Amount) / 1e8
			response.VolumeFee += float64(voucher.BitcouFeePaymentData.Amount) / 1e8
		}
		response.Total += 1
		switch voucher.Status {
		case hestia.GetVoucherStatusString(hestia.VoucherStatusPending):
			response.Pending += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusConfirming):
			response.Confirming += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusConfirmed):
			response.Confirmed += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusError):
			response.Error += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusAwaitingProvider):
			response.AwaitingProvider += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefundTotal):
			response.RefundTotal += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefundFee):
			response.RefundFee += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefunded):
			response.Refunded += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefundedPartially):
			response.RefundedPartially += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusComplete):
			response.Complete += 1
		}
	}
	return response, nil
}

func (sc StatsController) GetShiftsByTimeStats(userData hestia.User, params Params) (interface{}, error) {
	if params.Timestamp == "" {
		return nil, errors.New("please add the timestamp param")
	}
	shifts, err := sc.ShiftModel.GetAll("all")
	if err != nil {
		return nil, err
	}
	response := models.ShiftStatsResponse{
		Pending:    0,
		Confirming: 0,
		Confirmed:  0,
		Error:      0,
		Refund:     0,
		Refunded:   0,
		Complete:   0,
		Total:      0,
	}
	timestamp, err := strconv.ParseInt(params.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	for _, shift := range shifts {
		if shift.Timestamp < timestamp {
			continue
		}
		if shift.Status == hestia.GetShiftStatusString(hestia.ShiftStatusComplete) {
			req := obol.ObolRequest{ObolURL: "https://obol.polispay.com"}
			if shift.FeePayment.Txid != "" {
				feeRates, ok := coinRates[shift.FeePayment.Coin]
				if !ok {
					feeRates, err := req.GetCoinRates(shift.FeePayment.Coin)
					if err != nil {
						return nil, err
					}
					coinRates[shift.FeePayment.Coin] = feeRates
				}
				feeRates, err := req.GetCoinRates(shift.FeePayment.Coin)
				if err != nil {
					return nil, err
				}
				var btcFeeRate float64
				for _, rate := range feeRates {
					if rate.Code == "BTC" {
						btcFeeRate = rate.Rate
					}
				}
				response.Volume += btcFeeRate * (float64(shift.FeePayment.Amount) / 1e8)
			}
			paymentRates, ok := coinRates[shift.Payment.Coin]
			if !ok {
				paymentRates, err := req.GetCoinRates(shift.Payment.Coin)
				if err != nil {
					return nil, err
				}
				coinRates[shift.Payment.Coin] = paymentRates
			}

			var btcPaymentRate float64
			for _, rate := range paymentRates {
				if rate.Code == "BTC" {
					btcPaymentRate = rate.Rate
				}
			}
			response.Volume += btcPaymentRate * (float64(shift.Payment.Amount) / 1e8)
		}
		response.Total += 1
		switch shift.Status {
		case hestia.GetShiftStatusString(hestia.ShiftStatusPending):
			response.Pending += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusConfirming):
			response.Confirming += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusConfirmed):
			response.Confirmed += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusError):
			response.Error += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusRefund):
			response.Refund += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusRefunded):
			response.Refunded += 1
		case hestia.GetShiftStatusString(hestia.ShiftStatusComplete):
			response.Complete += 1
		}
	}
	return response, nil
}

func (sc StatsController) GetVouchersByTimeStats(userData hestia.User, params Params) (interface{}, error) {
	if params.Timestamp == "" {
		return nil, errors.New("please add the timestamp param")
	}
	vouchers, err := sc.VouchersModel.GetAll("all")
	if err != nil {
		return nil, err
	}
	response := models.VoucherStatsResponse{
		Pending:           0,
		Confirming:        0,
		Confirmed:         0,
		AwaitingProvider:  0,
		Error:             0,
		RefundFee:         0,
		RefundTotal:       0,
		RefundedPartially: 0,
		Refunded:          0,
		Complete:          0,
		Total:             0,
	}
	timestamp, err := strconv.ParseInt(params.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	for _, voucher := range vouchers {
		if voucher.Timestamp < timestamp {
			continue
		}
		if voucher.Status == hestia.GetVoucherStatusString(hestia.VoucherStatusComplete) {
			response.Volume += float64(voucher.BitcouPaymentData.Amount) / 1e8
			response.VolumeFee += float64(voucher.BitcouFeePaymentData.Amount / 1e8)
		}
		response.Total += 1
		switch voucher.Status {
		case hestia.GetVoucherStatusString(hestia.VoucherStatusPending):
			response.Pending += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusConfirming):
			response.Confirming += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusConfirmed):
			response.Confirmed += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusError):
			response.Error += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusAwaitingProvider):
			response.AwaitingProvider += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefundTotal):
			response.RefundTotal += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefundFee):
			response.RefundFee += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefunded):
			response.Refunded += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusRefundedPartially):
			response.RefundedPartially += 1
		case hestia.GetVoucherStatusString(hestia.VoucherStatusComplete):
			response.Complete += 1
		}
	}
	return response, nil
}
