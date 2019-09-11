package test

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/controllers"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
)

var (
	cardsCtrl    controllers.CardsController
	depositsCtrl controllers.DepositsController
	ordersCtrl   controllers.OrdersController
	shiftCtrl    controllers.ShiftsController
	userCtrl     controllers.UsersController
	vouchersCtrl controllers.VouchersController
	coinsCtrl    controllers.CoinsController
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
	vouchersModel := &models.VouchersModel{Db: db, Collection: "vouchers"}
	usersModel := &models.UsersModel{Db: db, Collection: "users"}
	coinsModel := &models.CoinsModel{Db: db, Collection: "coins"}

	// Init Controllers
	cardsCtrl = controllers.CardsController{Model: cardsModel, UserModel: usersModel}
	depositsCtrl = controllers.DepositsController{Model: depositsModel, UserModel: usersModel}
	ordersCtrl = controllers.OrdersController{Model: ordersModel, UserModel: usersModel}
	shiftCtrl = controllers.ShiftsController{Model: shiftsModel, UserModel: usersModel}
	userCtrl = controllers.UsersController{Model: usersModel}
	vouchersCtrl = controllers.VouchersController{Model: vouchersModel, UserModel: usersModel}
	coinsCtrl = controllers.CoinsController{Model: coinsModel}
}
