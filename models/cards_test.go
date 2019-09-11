package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
