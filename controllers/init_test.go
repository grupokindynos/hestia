package controllers

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
)

var (
	cardsCtrl    CardsController
	depositsCtrl DepositsController
	ordersCtrl   OrdersController
	shiftCtrl    ShiftsController
	coinsCtrl    CoinsController
)

func init() {
	_ = godotenv.Load("../.env")
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Init DB models
	shiftsModel := &models.ShiftModel{Db: db, Collection: "shifts"}
	cardsModel := &models.CardsModel{Db: db, Collection: "cards"}
	ordersModel := &models.OrdersModel{Db: db, Collection: "orders"}
	depositsModel := &models.DepositsModel{Db: db, Collection: "deposits"}
	usersModel := &models.UsersModel{Db: db, Collection: "users"}
	coinsModel := &models.CoinsModel{Db: db, Collection: "coins"}

	// Init Controllers
	cardsCtrl = CardsController{Model: cardsModel, UserModel: usersModel}
	depositsCtrl = DepositsController{Model: depositsModel, UserModel: usersModel}
	ordersCtrl = OrdersController{Model: ordersModel, UserModel: usersModel}
	shiftCtrl = ShiftsController{Model: shiftsModel, UserModel: usersModel}
	coinsCtrl = CoinsController{Model: coinsModel}
}
