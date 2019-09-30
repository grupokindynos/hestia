package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type CardsModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *CardsModel) Get(id string) (card hestia.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return card, err
	}
	err = doc.DataTo(&card)
	if err != nil {
		return card, err
	}
	return card, nil
}

func (m *CardsModel) Update(card hestia.Card) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(card.CardCode).Set(ctx, card)
	return err
}

func (m *CardsModel) GetAll() (cards []hestia.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var card hestia.Card
		_ = doc.DataTo(&card)
		cards = append(cards, card)
	}
	return cards, nil
}
