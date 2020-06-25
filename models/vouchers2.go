package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/herodotus"
	"github.com/grupokindynos/common/hestia"
	"strconv"
	"time"
)

type VouchersModelV2 struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *VouchersModelV2) Get(id string) (voucher hestia.VoucherV2, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return voucher, err
	}
	err = doc.DataTo(&voucher)
	if err != nil {
		return voucher, err
	}
	return voucher, nil
}

func (m *VouchersModelV2) Update(voucher hestia.VoucherV2) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(voucher.Id).Set(ctx, voucher)
	return err
}

func (m *VouchersModelV2) GetAll(filter int, timefilter string) (vouchers []hestia.VoucherV2, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	var docSnap []*firestore.DocumentSnapshot
	if timefilter != "" {
		timeInt, err := strconv.Atoi(timefilter)
		if err != nil {
			return nil, err
		}
		query := ref.Where("timestamp", ">=", timeInt)
		docSnap, err = query.Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		if filter == -1 {
			query := ref.OrderBy("timestamp", firestore.Asc)
			docSnap, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, err
			}
			/*docSnap, err = ref.Documents(ctx).GetAll()
			if err != nil {
				return nil, err
			}*/
		} else {
			query := ref.Where("status", "==", filter)
			docSnap, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, err
			}
		}
	}
	for _, doc := range docSnap {
		var voucher hestia.VoucherV2
		_ = doc.DataTo(&voucher)
		vouchers = append(vouchers, voucher)
	}
	return vouchers, nil
}

func (m *VouchersModelV2) GetWithComposedQuery(filters herodotus.VoucherV2Filters) (vouchers []hestia.VoucherV2, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	var docSnap []*firestore.DocumentSnapshot
	query := ref.Where("created_time", ">=", filters.FromTimestamp).Where("created_time", "<=", filters.ToTimestamp)

	docSnap, err = query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docSnap {
		var voucher hestia.VoucherV2
		_ = doc.DataTo(&voucher)
		if checkVoucherWithFilters(voucher, filters) {
			vouchers = append(vouchers, voucher)
		}
	}

	return vouchers, nil
}

func checkVoucherWithFilters(voucher hestia.VoucherV2, filters herodotus.VoucherV2Filters) bool {
	if len(filters.VoucherId) > 0 && !filters.VoucherId[voucher.VoucherId] {return false}
	if len(filters.UserId) > 0 && !filters.UserId[voucher.UserId] {return false}
	if len(filters.ProviderId) > 0 && !filters.ProviderId[voucher.ProviderId] {return false}
	if len(filters.Status) > 0 && !filters.Status[voucher.Status] {return false}
	if len(filters.Coin) > 0 && !filters.Coin[voucher.UserPayment.Coin] {return false}

	return true
}
