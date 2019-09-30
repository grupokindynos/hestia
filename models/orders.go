package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type OrdersModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *OrdersModel) Get(id string) (order hestia.Order, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return order, err
	}
	err = doc.DataTo(&order)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (m *OrdersModel) Update(order hestia.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(order.ID).Set(ctx, order)
	return err
}

func (m *OrdersModel) GetAll() (orders []hestia.Order, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var order hestia.Order
		_ = doc.DataTo(&order)
		orders = append(orders, order)
	}
	return orders, nil
}
