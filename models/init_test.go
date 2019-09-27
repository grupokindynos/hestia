package models

import (
	"github.com/grupokindynos/hestia/config"
	"github.com/joho/godotenv"
)

var (
	shiftsModel   *ShiftModel
	cardsModel    *CardsModel
	depositsModel *DepositsModel
	ordersModel   *OrdersModel
	vouchersModel *VouchersModel
	coinsModel    *CoinsModel
	configModel   *GlobalConfigModel
	usersModel    *UsersModel
)

func init() {
	_ = godotenv.Load("../.env")
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	// Init DB models
	shiftsModel = &ShiftModel{Db: db, Collection: "shifts"}
	cardsModel = &CardsModel{Db: db, Collection: "cards"}
	ordersModel = &OrdersModel{Db: db, Collection: "orders"}
	depositsModel = &DepositsModel{Db: db, Collection: "deposits"}
	vouchersModel = &VouchersModel{Db: db, Collection: "vouchers"}
	coinsModel = &CoinsModel{Db: db, Collection: "coins"}
	usersModel = &UsersModel{Db: db, Collection: "users"}
	configModel = &GlobalConfigModel{Db: db, Collection: "config"}
}
