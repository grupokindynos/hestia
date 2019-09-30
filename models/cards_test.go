package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCardsModel_Update(t *testing.T) {
	err := cardsModel.Update(TestCard)
	assert.Nil(t, err)
}

func TestCardsModel_Get(t *testing.T) {
	newCard, err := cardsModel.Get(TestCard.CardCode)
	assert.Nil(t, err)
	assert.Equal(t, TestCard, newCard)
}

func TestCardsModel_GetAll(t *testing.T) {
	cards, err := cardsModel.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(cards))
	assert.Equal(t, TestCard, cards[0])
	assert.IsType(t, []hestia.Card{}, cards)
}
