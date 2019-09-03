package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Card struct{}

type CardsModel struct {
	Db *mongo.Database
}

func (m *CardsModel) GetCard(cardcode string) (card Card, err error) {
	return card, err
}

func (m *CardsModel) UpdateCard(uid string, card Card) error {
	return nil
}

func (m *CardsModel) StoreCard(uid string, card Card) error {
	return nil
}
