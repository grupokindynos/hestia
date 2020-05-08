package models

import (
"cloud.google.com/go/firestore"
"context"
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

func (m *VouchersModelV2) GetAll(filter string, timefilter string) (vouchers []hestia.VoucherV2, err error) {
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
		if filter == "all" {
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
