package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var cardsCtrl CardsController

func init() {
	_ = godotenv.Load("../.env")
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	model := &models.CardsModel{
		Db:         db,
		Collection: "cards",
	}
	userModel := &models.UsersModel{
		Db:         db,
		Collection: "users",
	}
	cardsCtrl = CardsController{
		Model:     model,
		UserModel: userModel,
	}
}

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
