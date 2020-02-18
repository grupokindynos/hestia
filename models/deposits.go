package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type DepositsModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *DepositsModel) Get(id string) (deposit hestia.Deposit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return deposit, err
	}
	err = doc.DataTo(&deposit)
	if err != nil {
		return deposit, err
	}
	return deposit, nil
}

func (m *DepositsModel) Update(deposit hestia.Deposit) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(deposit.ID).Set(ctx, deposit)
	return err
}

func (m *DepositsModel) GetAll(filter string) (deposits []hestia.Deposit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	var docSnap []*firestore.DocumentSnapshot
	if filter == "all" {
		docSnap, err = ref.Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		query := ref.Where("status", "==", filter)
		docSnap, err = query.Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
	}
	for _, doc := range docSnap {
		var deposit hestia.Deposit
		_ = doc.DataTo(&deposit)
		deposits = append(deposits, deposit)
	}
	return deposits, nil
}
