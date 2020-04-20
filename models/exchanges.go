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

func (m *ExchangesModel) Get(id string) (exchange hestia.ExchangeInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return exchange, err
	}
	err = doc.DataTo(&exchange)
	if err != nil {
		return exchange, err
	}
	return exchange, nil
}

func (m *ExchangesModel) Update(exchange hestia.ExchangeInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(exchange.Id).Set(ctx, exchange)
	return err
}

func (m *ExchangesModel) GetAll() (exchanges []hestia.ExchangeInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var exchange hestia.ExchangeInfo
		_ = doc.DataTo(&exchange)
		exchanges = append(exchanges, exchange)
	}
	return exchanges, nil
}
