package models

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	shiftsModel    *ShiftModel
	cardsModel     *CardsModel
	depositsModel  *DepositsModel
	ordersModel    *OrdersModel
	vouchersModel  *VouchersModel
	coinsModel     *CoinsModel
	configModel    *GlobalConfigModel
	usersModel     *UsersModel
	exchangesModel *ExchangesModel
	balancesModel  *BalancesModel
)

func init() {
	_ = godotenv.Load("../.env")

	fbCredStr := os.Getenv("FIREBASE_CRED")
	fbCred, err := base64.StdEncoding.DecodeString(fbCredStr)
	if err != nil {
		log.Fatal("unable to decode firebase credential string:")
	}
	opt := option.WithCredentialsJSON(fbCred)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("unable to initialize firebase app")
	}

	// Init Database
	firestore, err := fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	baseDoc := firestore.Collection("polispay").Doc("hestia")

	// Init DB models
	shiftsModel = &ShiftModel{Firestore: baseDoc, Collection: "shifts"}
	cardsModel = &CardsModel{Firestore: baseDoc, Collection: "cards"}
	ordersModel = &OrdersModel{Firestore: baseDoc, Collection: "orders"}
	depositsModel = &DepositsModel{Firestore: baseDoc, Collection: "deposits"}
	vouchersModel = &VouchersModel{Firestore: baseDoc, Collection: "vouchers"}
	coinsModel = &CoinsModel{Firestore: baseDoc, Collection: "coins"}
	usersModel = &UsersModel{Firestore: baseDoc, Collection: "users"}
	configModel = &GlobalConfigModel{Firestore: baseDoc, Collection: "config"}
	exchangesModel = &ExchangesModel{Firestore: baseDoc, Collection: "exchanges"}
	balancesModel = &BalancesModel{Firestore: baseDoc, Collection: "balances"}
}
