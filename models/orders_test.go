package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestOrder = Order{
	ID:     "TEST-ORDER",
	UID:    "XYZ12345678910",
	Status: "COMPLETED",
	PaymentInfo: Payment{
		Address:       "FAKE-ADDRRESS",
		Amount:        "100",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	AddressInfo: AddressInformation{
		City:       "FAKE-CITY",
		Country:    "FAKE-COUNTRY",
		PostalCode: "00000",
		Email:      "TEST@TEST.COM",
		Street:     "FAKE-STREET",
	},
	Delivery: DeliveryOption{
		Type:    "DHL",
		Service: "EXPRESS",
		DeliveryAddress: AddressInformation{
			City:       "FAKE-CITY",
			Country:    "FAKE-COUNTRY",
			PostalCode: "00000",
			Email:      "TEST@TEST.COM",
			Street:     "FAKE-STREET",
		},
	},
	PersonalizationData: PersonalizationInformation{
		PersonalData: PersonalInformation{
			BirthDate:   "20190211",
			CivilState:  "MARRIED",
			FirstName:   "TEST-USER",
			LastName:    "TEST-USER",
			Sex:         "MALE",
			HomeNumber:  "00-00-00-00",
			Nationality: "MEXICAN",
			PassportID:  "1234567890",
		},
		AddressData: AddressInformation{
			City:       "FAKE-CITY",
			Country:    "FAKE-COUNTRY",
			PostalCode: "00000",
			Email:      "TEST@TEST.COM",
			Street:     "FAKE-STREET",
		},
	},
}

func TestOrdersModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	err = model.Update(TestOrder)
	assert.Nil(t, err)
}

func TestOrdersModel_Get(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	newOrder, err := model.Get(TestOrder.ID)
	assert.Nil(t, err)
	assert.Equal(t, TestOrder, newOrder)
}

func TestOrdersModel_GetAll(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := OrdersModel{
		Db:         db,
		Collection: "orders",
	}
	orders, err := model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(orders))
	assert.IsType(t, []Order{}, orders)
}
