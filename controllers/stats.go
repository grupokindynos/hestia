package controllers

import (
	"errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"strconv"
)

type StatsController struct {
	ShiftModel    *models.ShiftModel
	VouchersModel *models.VouchersModel
	DepositsModel *models.DepositsModel
	OrdersModel   *models.OrdersModel
}

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
