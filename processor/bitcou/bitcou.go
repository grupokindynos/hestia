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
	0:  "Wunschgutschein", // Missing img
	1:  "Access Wireless",
	2:  "Airvoice",
	3:  "AT&T",
	4:  "AT&T",
	5:  "AT&T",
	6:  "Bell",
	7:  "Black Wireless",
	8:  "ChatR",
	9:  "Cricket Wireless",
	10: "DollarPhone",
	11: "easyGO",
	12: "Feelsafe Wireless",
	13: "Fido",
	14: "Freedom Mobile",
	15: "Go", // Missing img
	16: "Google Play",
	17: "H2O",
	18: "i-Wireless", // Missing img
	19: "iTunes",
	20: "Koodo",
	21: "Life Wireless",
	22: "Lycamobile",
	23: "Metro", // Missing img
	24: "Movistar",
	25: "Bell MTS",
	26: "NET10",
	27: "Page Plus",
	28: "PlayStation",
	29: "Public Mobile",
	30: "Pure Unlimited",
	31: "Red Pocket Mobile",
	32: "Rogers",
	33: "Simple Mobile",
	34: "SIN PIN",
	35: "Skype",
	36: "Solo Mobile",
	37: "Southern Linc",
	38: "Spotify",
	39: "T-Mobile",
	40: "Telcel",
	41: "Telcel",
	42: "Total Wireless",
	43: "TracFone",
	44: "Ultra Mobile",
	45: "Unefon",
	46: "Verizon",
	47: "Virgin",
	48: "Xbox",
	49: "Xfinity",
	50: "Bildmobil", // Missing img
	51: "Blau", // Missing img
	52: "Confstar", // Missing img
	53: "E-Plus", // Missing img
	54: "FC Bayern Mobil", // Missing img
	55: "Fonic", // Missing img
	56: "Klarmobil", // Missing img
	57: "Lebara", // Missing img
	58: "Lifecell", // Missing img
	59: "Lycamobile",
	60: "Mobi", // Missing img
	61: "O2", // Missing img
	62: "Ortel", // Missing img
	63: "Otelo", // Missing img
	64: "SIM", // Missing img
	65: "Simyo", // Missing img
	66: "Tchibo Mobil", // Missing img
	67: "Telekom", // Missing img
	68: "Vodafone", // Missing img
	69: "Yourfone", // Missing img
	70: "About you", // Missing img
	71: "Adidas", // Missing img
	72: "Amazon", // Missing img
	73: "DAZN", // Missing img
	74: "Deezer", // Missing img
	75: "Epay Card", // Missing img
	76: "Google Play", // Missing img
	77: "iTunes",
	78: "Netflix", // Missing img
	79: "Spotify",
	80: "Zalando", // Missing img
	81: "Battle.net", // Missing img
	82: "Bigpoint", // Missing img
	83: "Nintendo", // Missing img
	84: "PlayStation",
	85: "Steam", // Missing img
	86: "Xbox",
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
