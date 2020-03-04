package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	obol "github.com/grupokindynos/hestia/processor/liquidity/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var CoinsModel *models.CoinsModel

const (
	MinVolumeForConversions = 5000
	MinVolumeForVouchers    = 2000
	MinVolumeForDeposits    = 10000
	MinVolumeForOrders      = 10000
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
	CoinsModel = &models.CoinsModel{Firestore: doc, Collection: "coins"}
}

func main() {
	// Load current coin configuration
	currentCoinConfigs, err := CoinsModel.GetCoinsData()
	if err != nil {
		log.Panic(err)
	}
	coinConfigMap := make(map[string]hestia.Coin)
	for _, coin := range currentCoinConfigs {
		coinConfigMap[coin.Ticker] = coin
	}
	var newCoinConfigs []hestia.Coin
	for _, coin := range coinfactory.Coins {
		// Omit BTC and POLIS because we are always available for this coins.
		if coin.Info.Tag == "BTC" || coin.Info.Tag == "POLIS" {
			continue
		}

		// Temporally omit ETH and ERC20 + ONION
		if coin.Info.Tag == "ETH" || coin.Info.Token || coin.Info.Tag == "ONION" {
			continue
		}

		adrestiaCoin := false
		// Coins available for adrestia
		if coin.Info.Tag == "DASH" {
			adrestiaCoin = true
		}

		coinLiquidity, err := getLiquidity(coin.Info.Tag)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(coinLiquidity, coin.Info.Tag)
		currentCoinInfo := coinConfigMap[coin.Info.Tag]
		var orderAvailable, depositAvailable, shiftAvailable, vouchersAvailable bool

		if coinLiquidity > MinVolumeForDeposits {
			depositAvailable = true
		} else {
			depositAvailable = false
		}
		if coinLiquidity > MinVolumeForVouchers {
			vouchersAvailable = true
		} else {
			vouchersAvailable = false
		}

		if coinLiquidity > MinVolumeForConversions {
			shiftAvailable = true
		} else {
			shiftAvailable = false
		}
		if coinLiquidity > MinVolumeForOrders {
			orderAvailable = true
		} else {
			orderAvailable = false
		}

		newConfig := hestia.Coin{
			Ticker: coin.Info.Tag,
			Shift: hestia.Properties{
				FeePercentage: currentCoinInfo.Shift.FeePercentage,
				Available:     shiftAvailable,
			},
			Deposits: hestia.Properties{
				FeePercentage: currentCoinInfo.Deposits.FeePercentage,
				Available:     depositAvailable,
			},
			Vouchers: hestia.Properties{
				FeePercentage: currentCoinInfo.Vouchers.FeePercentage,
				Available:     vouchersAvailable,
			},
			Orders: hestia.Properties{
				FeePercentage: currentCoinInfo.Orders.FeePercentage,
				Available:     orderAvailable,
			},
			Balances: currentCoinInfo.Balances,
			Adrestia: adrestiaCoin,
		}
		newCoinConfigs = append(newCoinConfigs, newConfig)
	}
	err = CoinsModel.UpdateCoinsData(newCoinConfigs)
	if err != nil {
		log.Panic(err)
	}
}

func getLiquidity(coin string) (float64, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get("https://obol.polispay.com/liquidity/" + coin)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	var response obol.Liquidity
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	return response.Data, nil
}
