package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	firebase "firebase.google.com/go"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/models"
	cmc "github.com/grupokindynos/hestia/processor/volume/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	list, err := getVolumeForCoins()
	var newCoinConfigs []hestia.Coin
	for _, coin := range coinfactory.Coins {
		if coin.Info.Tag == "POLIS" {
			continue
		}
		coinInfo, err := getCoinInfoFromList(list, coin.Info.Tag)
		if err != nil {
			log.Panic(err)
		}
		currentCoinInfo := coinConfigMap[coin.Info.Tag]
		var orderAvailable, depositAvailable, shiftAvailable, vouchersAvailable bool

		if !currentCoinInfo.Deposits.Available {
			depositAvailable = false
		} else {
			if *coinInfo.Quote.USD.Volume24h > MinVolumeForDeposits {
				depositAvailable = true
			} else {
				depositAvailable = false
			}
		}

		if !currentCoinInfo.Vouchers.Available {
			vouchersAvailable = false
		} else {
			if *coinInfo.Quote.USD.Volume24h > MinVolumeForVouchers {
				vouchersAvailable = true
			} else {
				vouchersAvailable = false
			}
		}

		if !currentCoinInfo.Shift.Available {
			shiftAvailable = false
		} else {
			if *coinInfo.Quote.USD.Volume24h > MinVolumeForConversions {
				shiftAvailable = true
			} else {
				shiftAvailable = false
			}
		}

		if !currentCoinInfo.Orders.Available {
			orderAvailable = false
		} else {
			if *coinInfo.Quote.USD.Volume24h > MinVolumeForOrders {
				orderAvailable = true
			} else {
				orderAvailable = false
			}
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
		}
		newCoinConfigs = append(newCoinConfigs, newConfig)
	}

	err = CoinsModel.UpdateCoinsData(newCoinConfigs)
	if err != nil {
		log.Panic(err)
	}
}

func getVolumeForCoins() (cmc.CoinsList, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		return cmc.CoinsList{}, err
	}
	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("CMC_API"))
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return cmc.CoinsList{}, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	var CoinsList cmc.CoinsList
	err = json.Unmarshal(respBody, &CoinsList)
	if err != nil {
		return cmc.CoinsList{}, err
	}
	return CoinsList, nil
}

func getCoinInfoFromList(list cmc.CoinsList, coinTag string) (cmc.CoinInfo, error) {
	for _, coin := range list.Data {
		if coin.Symbol == strings.ToUpper(coinTag) {
			return coin, nil
		}
	}
	return cmc.CoinInfo{}, nil
}
