package models

import (
	"go.mongodb.org/mongo-driver/mongo"
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
	Db *mongo.Database
}
