package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	firebase "firebase.google.com/go"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/plutus"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var BalanceModel *models.BalancesModel

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
}

func main() {
	coinsData, err := BalanceModel.GetBalances()
	if err != nil {
		log.Fatal(err)
	}
	coinsMap := make(map[string]hestia.CoinBalances)
	for _, coin := range coinsData {
		coinsMap[coin.Ticker] = coin
	}
	var newCoinData []hestia.CoinBalances
	for _, coin := range coinfactory.Coins {
		currCoinInfo := coinsMap[coin.Info.Tag]
		balanceReq, err := mvt.CreateMVTToken("GET", plutus.ProductionURL+"/balance/"+coin.Info.Tag, "hestia", os.Getenv("MASTER_PASSWORD"), nil, os.Getenv("PLUTUS_AUTH_USERNAME"), os.Getenv("PLUTUS_AUTH_PASSWORD"), os.Getenv("HESTIA_PRIVATE_KEY"))
		if err != nil {
			currCoinInfo.Balance = 0
			currCoinInfo.Status = "ERROR: Unable to create request"
			newCoinData = append(newCoinData, currCoinInfo)
			continue
		}
		client := http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 30,
		}
		res, err := client.Do(balanceReq)
		if err != nil {
			currCoinInfo.Balance = 0
			currCoinInfo.Status = "ERROR: Unable to do request"
			newCoinData = append(newCoinData, currCoinInfo)
			continue
		}
		if res.StatusCode == 200 {
			var resToken string
			respBody, _ := ioutil.ReadAll(res.Body)
			err = json.Unmarshal(respBody, &resToken)
			if err != nil {
				currCoinInfo.Balance = 0
				currCoinInfo.Status = "ERROR: Unable to unmarshal response body"
				newCoinData = append(newCoinData, currCoinInfo)
				continue
			}
			valid, payload := mrt.VerifyMRTToken(res.Header.Get("service"), resToken, os.Getenv("PLUTUS_PUBLIC_KEY"), os.Getenv("MASTER_PASSWORD"))
			if !valid {
				currCoinInfo.Balance = 0
				currCoinInfo.Status = "ERROR: Response us invalid"
				newCoinData = append(newCoinData, currCoinInfo)
				continue
			}
			var balance plutus.Balance
			err = json.Unmarshal(payload, &balance)
			if err != nil {
				currCoinInfo.Balance = 0
				currCoinInfo.Status = "ERROR: Unable to unmarshal response payload"
				newCoinData = append(newCoinData, currCoinInfo)
				continue
			}
			currCoinInfo.Balance = balance.Confirmed
			currCoinInfo.Status = "SUCCESS"
			newCoinData = append(newCoinData, currCoinInfo)
			continue
		} else {
			currCoinInfo.Balance = 0
			currCoinInfo.Status = "ERROR: Response Internal Error"
			newCoinData = append(newCoinData, currCoinInfo)
			continue
		}
	}
	err = BalanceModel.UpdateBalances(newCoinData)
	if err != nil {
		log.Panic(err)
	}
}
