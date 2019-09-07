package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Order struct {
	ID                  string                     `bson:"id" json:"id"`
	UID                 string                     `bson:"uid" json:"uid"`
	Status              string                     `bson:"status" json:"status"`
	PaymentInfo         Payment                    `bson:"payment_info" json:"payment_info"`
	AddressInfo         AddressInformation         `bson:"address_info" json:"address_info"`
	Delivery            DeliveryOption             `bson:"delivery" json:"delivery"`
	PersonalizationData PersonalizationInformation `bson:"personalization_data" json:"personalization_data"`
}

type PersonalizationInformation struct {
	PersonalData PersonalInformation `bson:"personal_data" json:"personal_data"`
	AddressData  AddressInformation  `bson:"address_data" json:"address_data"`
}

type PersonalInformation struct {
	BirthDate   string `bson:"birth_date" json:"birth_date"`
	CivilState  string `bson:"civil_state" json:"civil_state"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Sex         string `bson:"sex" json:"sex"`
	HomeNumber  string `bson:"home_number" json:"home_number"`
	Nationality string `bson:"nationality" json:"nationality"`
	PassportID  string `bson:"passport_id" json:"passport_id"`
}

type DeliveryOption struct {
	Type            string             `bson:"type" json:"type"`
	Service         string             `bson:"service" json:"service"`
	DeliveryAddress AddressInformation `bson:"delivery_address" json:"delivery_address"`
}

type AddressInformation struct {
	City       string `bson:"city" json:"city"`
	Country    string `bson:"country" json:"country"`
	PostalCode string `bson:"postal_code" json:"postal_code"`
	Email      string `bson:"email" json:"email"`
	Street     string `bson:"street" json:"street"`
}

type OrdersModel struct {
	Db *mongo.Database
}
