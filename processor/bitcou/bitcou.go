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
	1:  "Access Wireless",
	2:  "AirVoice Wireless",
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
	15: "Go Smart",
	16: "Google Play",
	17: "H2O",
	18: "i-Wireless",
	19: "iTunes",
	20: "Koodo",
	21: "Life Wireless",
	22: "Lycamobile",
	23: "T-Mobile Metro",
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
	50: "Bild Mobile",
	51: "Blau",
	52: "Congstar",
	53: "E-Plus",
	54: "FC Bayern Mobil",
	55: "Fonic",
	56: "Klarmobile",
	57: "Lebara",
	58: "Lifecell",
	59: "Lycamobile",
	60: "Mobi",
	61: "O2",
	62: "Ortel",
	63: "Otelo",
	64: "SIM",
	65: "Simyo",
	66: "Tchibo Mobile",
	67: "Telekom",
	68: "Vodafone",
	69: "Yourfone",
	70: "About You",
	71: "Adidas",
	72: "Amazon",
	73: "DAZN",
	74: "Deezer",
	75: "Epay Card",
	76: "Google Play",
	77: "iTunes",
	78: "Netflix",
	79: "Spotify",
	80: "Zalando",
	81: "Battle Net",
	82: "Big Point",
	83: "Nintendo",
	84: "PlayStation",
	85: "Steam",
	86: "Xbox",
	87: "Wunschgutschein",
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
			if availableCountry == "usa" &&
				voucher.ProviderID == 2 ||
				voucher.ProviderID == 3 ||
				voucher.ProviderID == 4 ||
				voucher.ProviderID == 5 ||
				voucher.ProviderID == 7 ||
				voucher.ProviderID == 15 ||
				voucher.ProviderID == 17 ||
				voucher.ProviderID == 21 ||
				voucher.ProviderID == 22 ||
				voucher.ProviderID == 59 ||
				voucher.ProviderID == 26 ||
				voucher.ProviderID == 27 ||
				voucher.ProviderID == 30 ||
				voucher.ProviderID == 40 ||
				voucher.ProviderID == 41 ||
				voucher.ProviderID == 43 ||
				voucher.ProviderID == 49 ||
				voucher.ProviderID == 10 ||
				voucher.ProviderID == 34  {
				continue
			}
			if voucher.ProviderID == 24 && voucher.Benefits["Mobile"] && voucher.Benefits["Minutes"] && voucher.Benefits["Data"] {
				voucher.ProviderName = "Movistar Bundles"
			} else if voucher.ProductID == 17 {
				voucher.ProviderName = "PlayStation Live"
			} else if voucher.ProductID == 361 {
				voucher.ProviderName = "Battle Net Warcraft"
			} else if voucher.ProductID == 363 {
				voucher.ProviderName = "Nintendo Switch"
			} else {
				voucher.ProviderName = provName
			}
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
