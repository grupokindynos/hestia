package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type ExchangesModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *ExchangesModel) Get(id string) (order hestia.AdrestiaOrder, err error) {
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

func (m *ExchangesModel) Update(order hestia.AdrestiaOrder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(order.ID).Set(ctx, order)
	return err
}

func (m *ExchangesModel) GetAll(includeComplete bool, sinceTimestamp int) (orders []hestia.AdrestiaOrder, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var order hestia.AdrestiaOrder
		_ = doc.DataTo(&order)
		if sinceTimestamp != 0 {
			if int(order.Time) >= sinceTimestamp {
				if includeComplete {
					orders = append(orders, order)
				} else {
					if order.Status != hestia.AdrestiaStatusCompleted {
						orders = append(orders, order)
					}
				}
			}
		} else {
			if includeComplete {
				orders = append(orders, order)
			} else {
				if order.Status != hestia.AdrestiaStatusCompleted {
					orders = append(orders, order)
				}
			}
		}
	}
	return orders, nil
}
