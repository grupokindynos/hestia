package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestCard = Card{
	Address:    "TEST-ADDRRESS",
	CardCode:   "TEST-CARDCODE",
	CardNumber: "123456778",
	City:       "TEST-CITY",
	Email:      "TEST@TEST.COM",
	FirstName:  "TEST",
	LastName:   "CARD",
	UID:        "XYZ12345678910",
}

func TestCardsModel_Update(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := CardsModel{
		Db:         db,
		Collection: "cards",
	}
	err = model.Update(TestCard)
	assert.Nil(t, err)
}

func TestCardsModel_Get(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := CardsModel{
		Db:         db,
		Collection: "cards",
	}
	newCard, err := model.Get(TestCard.CardCode)
	assert.Nil(t, err)
	assert.Equal(t, TestCard, newCard)
}

func TestCardsModel_GetAll(t *testing.T) {
	db, err := config.ConnectDB()
	assert.Nil(t, err)
	model := CardsModel{
		Db:         db,
		Collection: "cards",
	}
	cards, err := model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(cards))
	assert.IsType(t, []Card{}, cards)
}
