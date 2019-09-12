package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCardsController_GetUserAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	cards, err := cardsCtrl.GetUserAll(models.TestUser, c)
	assert.Nil(t, err)
	var cardsArray []models.Card
	cardBytes, err := json.Marshal(cards)
	assert.Nil(t, err)
	err = json.Unmarshal(cardBytes, &cardsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Card{}, cards)
	assert.Equal(t, models.TestCard, cardsArray[0])
}

func TestCardsController_GetUserSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "cardcode", Value: models.TestCard.CardCode}}
	card, err := cardsCtrl.GetUserSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Card{}, card)
	assert.Equal(t, models.TestCard, card)
}

func TestCardsController_GetAll(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	cards, err := cardsCtrl.GetAll(models.TestUser, c)
	assert.Nil(t, err)
	var cardsArray []models.Card
	cardBytes, err := json.Marshal(cards)
	assert.Nil(t, err)
	err = json.Unmarshal(cardBytes, &cardsArray)
	assert.Nil(t, err)
	assert.IsType(t, []models.Card{}, cards)
	assert.Equal(t, models.TestCard, cardsArray[0])
}

func TestCardsController_GetSingle(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "cardcode", Value: models.TestCard.CardCode}}
	card, err := cardsCtrl.GetSingle(models.TestUser, c)
	assert.Nil(t, err)
	assert.IsType(t, models.Card{}, card)
	assert.Equal(t, models.TestCard, card)
}
