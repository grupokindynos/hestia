package controllers

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	cardsCtrl    CardsController
	depositsCtrl DepositsController
	ordersCtrl   OrdersController
	shiftCtrl    ShiftsController
	coinsCtrl    CoinsController
	globalCtrl   GlobalConfigController
	vouchersCtrl VouchersController
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
	baseDoc := firestore.Collection("polispay").Doc("hestia_test")

	// Init DB models
	shiftsModel := &models.ShiftModel{Firestore: baseDoc, Collection: "shifts"}
	cardsModel := &models.CardsModel{Firestore: baseDoc, Collection: "cards"}
	ordersModel := &models.OrdersModel{Firestore: baseDoc, Collection: "orders"}
	depositsModel := &models.DepositsModel{Firestore: baseDoc, Collection: "deposits"}
	usersModel := &models.UsersModel{Firestore: baseDoc, Collection: "users"}
	coinsModel := &models.CoinsModel{Firestore: baseDoc, Collection: "coins"}
	globalModel := &models.GlobalConfigModel{Firestore: baseDoc, Collection: "config"}
	vouchersModel := &models.VouchersModel{Firestore: baseDoc, Collection: "vouchers"}

	// Init Controllers
	cardsCtrl = CardsController{Model: cardsModel, UserModel: usersModel}
	depositsCtrl = DepositsController{Model: depositsModel, UserModel: usersModel}
	ordersCtrl = OrdersController{Model: ordersModel, UserModel: usersModel}
	shiftCtrl = ShiftsController{Model: shiftsModel, UserModel: usersModel}
	coinsCtrl = CoinsController{Model: coinsModel}
	globalCtrl = GlobalConfigController{Model: globalModel}
	vouchersCtrl = VouchersController{Model: vouchersModel, UserModel: usersModel}

}
