package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services/bitcou"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

// This tool must be run every 12 hours to index the bitcou vouchers list.

var ProvidersMap = map[int]string{
	0:  "Other",
	1:  "Telcel",
	2:  "AT&T",
	3:  "Unefon",
	4:  "Black Wireless",
	5:  "Black Wireless",
	6:  "Ultra Mobile PayGo",
	7:  "H2O Wireless",
	8:  "AT&T",
	9:  "H2O Wireless",
	10: "Access Wireless",
	11: "easyGO",
	12: "Verizon Wireless",
	13: "Life Wireless",
	14: "H2O Wireless",
	15: "H2O Wireless",
	16: "H2O Wireless",
	17: "H2O Wireless",
	18: "H2O Wireless",
	19: "Red Pocket",
	20: "Airvoice",
	21: "Airvoice",
	22: "Pure Unlimited",
	23: "Airvoice",
	24: "Pure Unlimited",
	25: "Pure Unlimited",
	26: "Pure",
	27: "i-Wireless Kroger",
	28: "Feelsafe",
	29: "Black Wireless",
	30: "AT&T",
	31: "Southern Linc",
	32: "DollarPhone",
	33: "Pure Unlimited",
	34: "Black Wireless",
	35: "AT&T",
	36: "Airvoice",
	37: "Life Wireless",
	38: "Life Wireless",
	39: "NET10 Wireless",
	40: "AT&T",
	41: "SIN PIN",
	42: "T-Mobile",
	43: "Page Plus",
	44: "TracFone",
	45: "Simple Mobile",
	46: "Telcel",
	47: "NET10 Wireless",
	48: "NET10 Wireless",
	49: "NET10 Wireless",
	50: "NET10 Wireless",
	51: "NET10 Wireless",
	53: "Lycamobile",
	54: "Go Smart",
	55: "iTunes",
	56: "Google Play",
	57: "Playstation Network",
	58: "Xbox",
	59: "T-Mobile",
	60: "Ultra Mobile",
	61: "T-Mobile",
	62: "Cricket Paygo",
	63: "Xfinity Prepaid Internet",
	64: "Xfinity Prepaid TV",
	65: "Xfinity Prepaid TV",
	66: "Total Wireless",
	67: "Bell",
	68: "ChatR",
	69: "Fido",
	70: "Koodo",
	71: "Bell MTS",
	72: "Public Mobile",
	73: "Rogers",
	74: "Virgin",
	75: "Freedom Mobile",
	76: "Solo Mobile",
	77: "Movistar",
	78: "Simple Mobile",
	79: "Skype",
	80: "Movistar",
	81: "Movistar Bundles",
	82: "Spotify",
	84: "Telcel",
	85: "Telcel",
}

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
	doc := firestore.Collection("bitcou")
	model := models.BitcouModel{Firestore: doc}
	service := bitcou.InitService()
	voucherList, err := service.GetList()
	if err != nil {
		panic("unable to load bitcou voucher list")
	}
	var countries []models.BitcouCountry
	var availableCountry []string
	for key, _ := range voucherList[0].Countries {
		availableCountry = append(availableCountry, key)
	}
	for _, availableCountry := range availableCountry {
		newCountryData := models.BitcouCountry{
			ID:       availableCountry,
			Vouchers: []bitcou.Voucher{},
		}
		for _, voucher := range voucherList {
			provName, ok := ProvidersMap[voucher.ProviderID]
			if !ok {
				continue
			}
			voucher.ProviderName = provName
			available := voucher.Countries[availableCountry]
			if available {
				newCountryData.Vouchers = append(newCountryData.Vouchers, voucher)
			}
		}
		countries = append(countries, newCountryData)
	}

	for _, bitcouCountry := range countries {
		err = model.AddCountry(bitcouCountry)
		if err != nil {
			panic("unable to store country information")
		}
	}

}
