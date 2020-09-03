package main

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"firebase.google.com/go"
	"github.com/grupokindynos/common/aes"
	"github.com/grupokindynos/common/hestia"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	hestiaDoc        *firestore.DocumentRef
	storageBuck      *storage.BucketHandle
	definedTimeLapse string
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
	// Init Bucket
	storageClient, err := fbApp.Storage(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	definedTimeLapse = os.Getenv("TIMELAPSE_CRON")
	if definedTimeLapse == "" {
		log.Fatal("not defined time lapse")
	}
	selectedBucket := "polispay-backups_"
	switch definedTimeLapse {
	case "hourly":
		selectedBucket += "hourly"
	case "daily":
		selectedBucket += "daily"
	case "weekly":
		selectedBucket += "weekly"
	case "monthly":
		selectedBucket += "monthly"
	default:
		log.Fatal("time lapse doesn't match")
	}
	storageBuck, err = storageClient.Bucket(selectedBucket)
	if err != nil {
		log.Fatal(err)
	}
	// Init Database
	client, err := fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	hestiaDoc = client.Collection("polispay").Doc("hestia")
}

/*

	The script purpose is to create a snapshot of the entire PolisPay DB under the /polispay collection into a
	known Golang Model and store it.

	Hestia database structure
	/polispay/hestia -> For production
	/polispay/hestia-test -> For unit testing.

	balances -> Ticker / hestia.CoinBalances | []hestia.CoinBalances
	cards -> CardCode / hestia.Cards | []hestia.Cards
	coins -> Ticker / hestia.Coin | []hestia.Coin
	config -> Service / hestia.Config | hestia.Config
	deposits -> DepositID / hestia.Deposit | []hestia.Deposit
    exchanges -> OrderID / hestia.AdrestiaOrder | []hestia.AdrestiaOrder
    orders -> OrderID / hestia.Order | []hestia.Order
    shifts -> ShiftID / hestia.Shift | []hestia.Shift
    users -> UID / hestia.User | []hestia.User
    vouchers -> VoucherID / hestia.Voucher | []hestia.Voucher

*/

func main() {
	balances, err := getBalances()
	if err != nil {
		log.Fatal(errors.New("unable to load balances"))
	}
	cards, err := getCards()
	if err != nil {
		log.Fatal(errors.New("unable to load cards"))
	}
	coins, err := getCoins()
	if err != nil {
		log.Fatal(errors.New("unable to load coins"))
	}
	config, err := getConfig()
	if err != nil {
		log.Fatal(errors.New("unable to load config"))
	}
	deposits, err := getDeposits()
	if err != nil {
		log.Fatal(errors.New("unable to load exchange deposits"))
	}
	exchanges, err := getExchanges()
	if err != nil {
		log.Fatal(errors.New("unable to load adrestia orders"))
	}
	orders, err := getOrders()
	if err != nil {
		log.Fatal(errors.New("unable to load orders"))
	}
	shifts, err := getShifts()
	if err != nil {
		log.Fatal(errors.New("unable to load shifts"))
	}
	users, err := getUsers()
	if err != nil {
		log.Fatal(errors.New("unable to load users"))
	}
	vouchers, err := getVouchers()
	if err != nil {
		log.Fatal(errors.New("unable to load vouchers"))
	}
	withdrawals, err := getWithdrawals()
	if err != nil {
		log.Fatal(errors.New("unable to load withdrawals"))
	}
	exchangeDeposits, err := getExchangeDeposits()
	if err != nil {
		log.Fatal(errors.New("unable to load deposits"))
	}
	balancerOrders, err := getBalancerOrders()
	if err != nil {
		log.Fatal(errors.New("unable to load balancerOrders"))
	}
	balancers, err := getBalancers()
	if err != nil {
		log.Fatal(errors.New("unable to load balancers"))
	}

	fullDB := hestia.HestiaDB{
		Balances:            balances,
		Cards:               cards,
		Coins:               coins,
		Config:              config,
		Deposits:            deposits,
		Exchanges:           exchanges,
		Orders:              orders,
		Shifts:              shifts,
		Users:               users,
		Vouchers:            vouchers,
		AdrestiaBalancer:    balancers,
		AdrestiaDeposits:    exchangeDeposits,
		AdrestiaOrders:      balancerOrders,
		AdrestiaWithdrawals: withdrawals,
	}
	jsonObj, err := json.Marshal(fullDB)
	if err != nil {
		log.Fatal(err)
	}
	encoded := base64.StdEncoding.EncodeToString(jsonObj)
	encrypted, err := aes.Encrypt([]byte(os.Getenv("HESTIA_BK_ENCRYPTION_KEY")), []byte(encoded))
	if err != nil {
		log.Fatal(err)
	}
	name := "polispay-database_"
	switch definedTimeLapse {
	case "hourly":
		currTime := time.Now()
		hour, min, _ := currTime.Clock()
		name += strconv.Itoa(hour) + ":" + strconv.Itoa(min)
	case "daily":
		currTime := time.Now()
		_, month, day := currTime.Date()
		name += month.String() + "-" + strconv.Itoa(day)
	case "weekly":
		currTime := time.Now()
		_, month, day := currTime.Date()
		name += month.String() + "-" + strconv.Itoa(day)
	case "monthly":
		currTime := time.Now()
		year, month, _ := currTime.Date()
		name += month.String() + "-" + strconv.Itoa(year)
	}
	w := storageBuck.Object(name).NewWriter(context.Background())
	w.ContentType = "text/plain"
	_, err = w.Write([]byte(encrypted))
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func getBalances() ([]hestia.CoinBalances, error) {
	collection := hestiaDoc.Collection("balances")
	docIter := collection.Documents(context.Background())
	var array []hestia.CoinBalances
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.CoinBalances
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getCards() ([]hestia.Card, error) {
	collection := hestiaDoc.Collection("cards")
	docIter := collection.Documents(context.Background())
	var array []hestia.Card
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Card
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getCoins() ([]hestia.Coin, error) {
	collection := hestiaDoc.Collection("coins")
	docIter := collection.Documents(context.Background())
	var array []hestia.Coin
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Coin
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getConfig() (hestia.Config, error) {
	return hestia.Config{}, nil
}

func getDeposits() ([]hestia.Deposit, error) {
	collection := hestiaDoc.Collection("deposits")
	docIter := collection.Documents(context.Background())
	var array []hestia.Deposit
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Deposit
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getExchanges() ([]hestia.ExchangeInfo, error) {
	collection := hestiaDoc.Collection("exchanges")
	docIter := collection.Documents(context.Background())
	var array []hestia.ExchangeInfo
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.ExchangeInfo
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getOrders() ([]hestia.Order, error) {
	collection := hestiaDoc.Collection("orders")
	docIter := collection.Documents(context.Background())
	var array []hestia.Order
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Order
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getShifts() ([]hestia.Shift, error) {
	collection := hestiaDoc.Collection("shifts")
	docIter := collection.Documents(context.Background())
	var array []hestia.Shift
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Shift
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getUsers() ([]hestia.User, error) {
	collection := hestiaDoc.Collection("users")
	docIter := collection.Documents(context.Background())
	var array []hestia.User
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.User
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getVouchers() ([]hestia.Voucher, error) {
	collection := hestiaDoc.Collection("vouchers")
	docIter := collection.Documents(context.Background())
	var array []hestia.Voucher
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Voucher
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getWithdrawals() ([]hestia.SimpleTx, error) {
	collection := hestiaDoc.Collection("adrestia_withdrawals")
	docIter := collection.Documents(context.Background())
	var array []hestia.SimpleTx
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.SimpleTx
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getExchangeDeposits() ([]hestia.SimpleTx, error) {
	collection := hestiaDoc.Collection("adrestia_deposits")
	docIter := collection.Documents(context.Background())
	var array []hestia.SimpleTx
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.SimpleTx
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getBalancerOrders() ([]hestia.BalancerOrder, error) {
	collection := hestiaDoc.Collection("adrestia_orders")
	docIter := collection.Documents(context.Background())
	var array []hestia.BalancerOrder
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.BalancerOrder
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}

func getBalancers() ([]hestia.Balancer, error) {
	collection := hestiaDoc.Collection("adrestia_balancer")
	docIter := collection.Documents(context.Background())
	var array []hestia.Balancer
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var elem hestia.Balancer
		err = doc.DataTo(&elem)
		if err != nil {
			return nil, err
		}
		array = append(array, elem)
	}
	return array, nil
}
