package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCardsModel_Update(t *testing.T) {
	err := cardsCtrl.Model.Update(TestCard)
	assert.Nil(t, err)
}

func TestCardsModel_Get(t *testing.T) {
	newCard, err := cardsCtrl.Model.Get(TestCard.CardCode)
	assert.Nil(t, err)
	assert.Equal(t, TestCard, newCard)
}

func TestCardsModel_GetAll(t *testing.T) {
	cards, err := cardsCtrl.Model.GetAll()
	assert.Nil(t, err)
	assert.NotZero(t, len(cards))
	assert.IsType(t, []models.Card{}, cards)
}

func TestCardsController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	cards, err := cardsCtrl.GetUserAll(TestUser, c)
	assert.Nil(t, err)
	var cardsArray []models.Card
	cardBytes, err := json.Marshal(cards)
	assert.Nil(t, err)
	err = json.Unmarshal(cardBytes, &cardsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Card{}, cards)
	assert.Equal(t, TestCard, cardsArray[0])
}

func TestCardsController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "cardcode", Value: TestCard.CardCode}}
	card, err := cardsCtrl.GetUserSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Card{}, card)
	assert.Equal(t, TestCard, card)
}

func TestCardsController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	cards, err := cardsCtrl.GetAll(TestUser, c)
	assert.Nil(t, err)
	var cardsArray []models.Card
	cardBytes, err := json.Marshal(cards)
	assert.Nil(t, err)
	err = json.Unmarshal(cardBytes, &cardsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Card{}, cards)
	assert.Equal(t, TestCard, cardsArray[0])
}

func TestCardsController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "cardcode", Value: TestCard.CardCode}}
	card, err := cardsCtrl.GetSingle(TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Card{}, card)
	assert.Equal(t, TestCard, card)
}
