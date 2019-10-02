package controllers

import (
	"encoding/json"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"testing"
)



func TestCardsController_GetUserAll(t *testing.T) {
	cards, err := cardsCtrl.GetAll(models.TestUser, TestParams)
	assert.Nil(t, err)
	var cardsArray []hestia.Card
	cardBytes, err := json.Marshal(cards)
	assert.Nil(t, err)
	err = json.Unmarshal(cardBytes, &cardsArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Card{}, cards)
	assert.Equal(t, models.TestCard, cardsArray[0])
}

func TestCardsController_GetUserSingle(t *testing.T) {
	card, err := cardsCtrl.GetSingle(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Card{}, card)
	assert.Equal(t, models.TestCard, card)
}

func TestCardsController_GetAll(t *testing.T) {
	cards, err := cardsCtrl.GetAll(models.TestUser, TestParams)
	assert.Nil(t, err)
	var cardsArray []hestia.Card
	cardBytes, err := json.Marshal(cards)
	assert.Nil(t, err)
	err = json.Unmarshal(cardBytes, &cardsArray)
	assert.Nil(t, err)
	assert.IsType(t, []hestia.Card{}, cards)
	assert.Equal(t, models.TestCard, cardsArray[0])
}

func TestCardsController_GetSingle(t *testing.T) {
	card, err := cardsCtrl.GetSingle(models.TestUser, TestParams)
	assert.Nil(t, err)
	assert.IsType(t, hestia.Card{}, card)
	assert.Equal(t, models.TestCard, card)
}
