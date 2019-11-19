package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services/bitcou"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

// This tool must be run every 12 hours to index the bitcou vouchers list.

func main() {
	_ = godotenv.Load()
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
	doc := firestore.Collection("polispay").Doc("bitcou")
	model := models.BitcouModel{Firestore: doc}
	service := bitcou.InitService()
	voucherList, err := service.GetList()
	if err != nil {
		panic("unable to load bitcou voucher list")
	}
	var availableCountries models.BitcouCountries
	for name, _ := range voucherList[0].Countries {
		availableCountries.Countries = append(availableCountries.Countries, name)
	}
	err = model.StoreCountries(availableCountries)
	if err != nil {
		fmt.Println(err)
		panic("unable to store bitcou available countries")
	}
}
