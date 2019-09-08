package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Card struct {
	Address    string `bson:"address" json:"address"`
	CardCode   string `bson:"card_code" json:"card_code"`
	CardNumber string `bson:"card_number" json:"card_number"`
	City       string `bson:"city" json:"city"`
	Email      string `bson:"email" json:"email"`
	FirstName  string `bson:"firstname" json:"firstname"`
	LastName   string `bson:"lastname" json:"lastname"`
	UID        string `bson:"uid" json:"uid"`
}

type Pin struct {
	CardCode string `bson:"card_code" json:"card_code"`
	PinCode  string `bson:"pin_code" json:"pin_code"`
}

type CardsModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *CardsModel) Get(id string) (card Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&card)
	return card, err
}

func (m *CardsModel) Update(id string, card Card) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: card}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *CardsModel) GetAll() (cards []Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, err := col.Find(ctx, nil)
	if err != nil {
		return cards, err
	}
	for curr.Next(ctx) {
		var card Card
		err := curr.Decode(&card)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}
	return cards, err
}
