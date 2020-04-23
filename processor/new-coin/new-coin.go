package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"flag"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

var (
	BalanceModel *models.BalancesModel
	CoinsModel   *models.CoinsModel
)

func init() {
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
	doc := firestore.Collection("polispay").Doc("hestia")
	BalanceModel = &models.BalancesModel{Firestore: doc, Collection: "balances"}
	CoinsModel = &models.CoinsModel{Firestore: doc, Collection: "coins"}
}

func main() {
	var ticker string
	flag.StringVar(&ticker, "ticker", "", "coin ticker to add")
	flag.Parse()
	if ticker == "" {
		panic("no ticker defined")
	}
	ticker = strings.ToUpper(ticker)
	newBalance := hestia.CoinBalances{
		Ticker:  ticker,
		Balance: 0,
		Status:  "SUCCESS",
	}
	err := BalanceModel.AddBalance(newBalance)
	if err != nil {
		panic(err)
	}
	newCoinModel := hestia.Coin{
		Ticker: ticker,
		Shift: hestia.Properties{
			FeePercentage: 1,
			Available:     false,
		},
		Deposits: hestia.Properties{
			FeePercentage: 1,
			Available:     false,
		},
		Vouchers: hestia.Properties{
			FeePercentage: 4.5,
			Available:     false,
		},
		Orders: hestia.Properties{
			FeePercentage: 1,
			Available:     false,
		},
		Balances: hestia.BalanceLimits{
			HotWallet: 0,
			Exchanges: 0,
		},
		Adrestia: hestia.AdrestiaInfo{
			Available: false,
			CoinUsage: 0,
		},
	}
	err = CoinsModel.AddCoin(newCoinModel)
	if err != nil {
		panic(err)
	}
}
