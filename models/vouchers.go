package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type VouchersModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *VouchersModel) Get(id string) (voucher hestia.Voucher, err error) {
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

func (m *VouchersModel) Update(voucher hestia.Voucher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(voucher.ID).Set(ctx, voucher)
	return err
}

func (m *VouchersModel) GetAll(filter string) (vouchers []hestia.Voucher, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var voucher hestia.Voucher
		_ = doc.DataTo(&voucher)
		if filter == "all" {
			vouchers = append(vouchers, voucher)
		} else {
			if voucher.Status == filter {
				vouchers = append(vouchers, voucher)
			}
		}
	}
	return vouchers, nil
}
