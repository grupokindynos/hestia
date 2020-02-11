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
var ProvidersMap = map[int]string{
	1 : "Access Wireless",
	2 : "Airvoice Wireless",
	4 : "AT&T",
	5 : "AT&T/Iusacell",
	6 : "Bell",
	7 : "Black Wireless",
	8 : "ChatR",
	9 : "Cricket Wireless",
	11 : "easyGO",
	12 : "Feelsafe Wireless",
	13 : "Fido",
	14 : "Freedom Mobile",
	15 : "Go",
	16 : "Google",
	17 : "H2O",
	18 : "i-Wireless",
	19 : "iTunes",
	20 : "Koodo",
	21 : "Life Wireless",
	22 : "Lycamobile",
	23 : "T-Mobile Metro",
	24 : "Movistar",
	25 : "Bell MTS",
	26 : "NET10",
	27 : "Page Plus",
	28 : "Playstation",
	29 : "Public Mobile",
	30 : "Pure Unlimited",
	31 : "Red Pocket Mobile",
	32 : "Rogers",
	33 : "Simple Mobile",
	35 : "Skype",
	36 : "Solo Mobile",
	37 : "Southern Linc",
	38 : "Spotify",
	39 : "T-Mobile",
	40 : "Telcel",
	41 : "Telcel",
	42 : "Total Wireless",
	43 : "TracFone",
	44 : "Ultra Mobile",
	45 : "Unefon",
	46 : "Verizon",
	47 : "Virgin",
	48 : "Xbox",
	49 : "Xfinity",
	50 : "Bildmobil",
	51 : "blau",
	52 : "congstar",
	53 : "E-Plus",
	54 : "FCBayernMobil",
	55 : "Fonic",
	56 : "Klarmobil",
	57 : "Lebara",
	58 : "Lifecell",
	59 : "Lycamobil",
	60 : "Mobi",
	61 : "O2",
	62 : "Ortel",
	63 : "Otelo",
	64 : "SIM",
	65 : "Simyo",
	66 : "TchiboMobil",
	67 : "Telekom",
	68 : "Vodafone",
	69 : "yourfone",
	70 : "aboutyoude",
	71 : "adidasde",
	72 : "Amazon",
	73 : "DAZNde",
	74 : "Deezerde",
	75 : "epaycard",
	76 : "GooglePlay",
	77 : "iTunesde",
	78 : "NetflixEUR",
	79 : "Spotifyde",
	80 : "zalandode",
	81 : "Battle.netEUR",
	82 : "Bigpoint",
	83 : "Nintendo",
	84 : "PlayStation",
	85 : "SteamEUR",
	86 : "Xbox",
	87 : "Wunschgutschein",
	89 : "Amazon",
	90 : "Best",
	91 : "Boston",
	92 : "Claro",
	94 : "Kolby",
	95 : "Netflix",
	96 : "Nintendo",
	97 : "Outback",
	98 : "Starbucks",
	99 : "TheHomeDepot",
}

var devProvidersMap = map[int]string{
	89 : "about",
	90 : "Access",
	91 : "adidas",
	92 : "Airvoice",
	93 : "Alo",
	94 : "Amazon",
	96 : "AT&T",
	97 : "AT&T/Iusacell",
	98 : "Battle.net",
	99 : "Bell",
	100 : "Bigpoint",
	101 : "Bildmobil",
	102 : "Black",
	103 : "blau",
	104 : "ChatR",
	105 : "congstar",
	106 : "Cricket",
	107 : "DAZN",
	108 : "Deezer",
	110 : "E-Plus",
	111 : "easyGO",
	112 : "epay",
	113 : "FC",
	114 : "Feelsafe",
	115 : "Fido",
	116 : "Fonic",
	117 : "Freedom",
	118 : "Go",
	119 : "Google",
	120 : "Google",
	121 : "H2O",
	122 : "i-Wireless",
	123 : "iTunes",
	124 : "iTunes",
	125 : "Klarmobil",
	126 : "Koodo",
	127 : "Lebara",
	128 : "Life",
	129 : "Lifecell",
	130 : "Lycamobil",
	131 : "Lycamobile",
	132 : "Metro",
	133 : "Mobi",
	134 : "Movistar",
	135 : "MTS",
	136 : "NET10",
	137 : "Netflix",
	138 : "Nintendo",
	139 : "O2",
	140 : "Ortel",
	141 : "Otelo",
	142 : "Page",
	143 : "PlayStation",
	144 : "Playstation",
	145 : "Public",
	146 : "Pure",
	147 : "Red",
	148 : "Rogers",
	149 : "SIM",
	150 : "Simple",
	151 : "Simyo",
	153 : "Skype",
	154 : "Solo",
	155 : "Southern",
	156 : "Spotify",
	157 : "Spotify",
	158 : "Steam",
	159 : "T-Mobile",
	160 : "Tchibo",
	161 : "Telcel",
	162 : "Telcel",
	163 : "Telekom",
	164 : "Total",
	165 : "TracFone",
	166 : "Ultra",
	167 : "Unefon",
	168 : "Verizon",
	169 : "Virgin",
	170 : "Vodafone",
	171 : "Wunschgutschein",
	172 : "Xbox",
	173 : "Xbox",
	174 : "Xfinity",
	175 : "yourfone",
	176 : "zalando",
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
	docTest := firestore.Collection("bitcou_test")
	model := models.BitcouModel{Firestore: doc, FirestoreTest: docTest}
	service := bitcou.InitService()
	voucherListProd, err := service.GetList(false)
	if err != nil {
		panic("unable to load bitcou voucher list")
	}
	voucherListDev, err := service.GetList(true)
	if err != nil {
		panic("unable to load bitcou voucher list")
	}
	var countries []models.BitcouCountry
	var countriesDev []models.BitcouCountry
	var availableCountry []string
	for key, _ := range voucherListProd[0].Countries {
		availableCountry = append(availableCountry, key)
	}
	for _, availableCountry := range availableCountry {
		newCountryData := models.BitcouCountry{
			ID:       availableCountry,
			Vouchers: []bitcou.Voucher{},
		}
		for _, voucher := range voucherListDev {
			fmt.Println(voucher.Name, " ", voucher.ProductID, " ", voucher.ProductID, " ", voucher.ProviderName)
			provName, ok := devProvidersMap[voucher.ProviderID]
			if !ok {
				continue
			}
			if voucher.ProviderID == 24 && voucher.Benefits["Mobile"] && voucher.Benefits["Minutes"] && voucher.Benefits["Data"] {
				fmt.Println("249")
				voucher.ProviderName = "Movistar Bundles"
			} else if voucher.ProductID == 17 {
				fmt.Println("252")
				voucher.ProviderName = "PlayStation Live"
			} else if voucher.ProductID == 361 {
				fmt.Println("255")
				voucher.ProviderName = "Battle Net Warcraft"
			} else if voucher.ProductID == 363 {
				fmt.Println("258")
				voucher.ProviderName = "Nintendo Switch"
			} else {
				fmt.Println("261")
				voucher.ProviderName = provName
			}
			fmt.Println("264")
			available := voucher.Countries[availableCountry]
			if available {
				fmt.Println("267")
				newCountryData.Vouchers = append(newCountryData.Vouchers, voucher)
			}
		}
		countriesDev = append(countriesDev, newCountryData)
	}
	fmt.Println("273")
	for _, availableCountry := range availableCountry {
		newCountryData := models.BitcouCountry{
			ID:       availableCountry,
			Vouchers: []bitcou.Voucher{},
		}
		for _, voucher := range voucherListProd {
			provName, ok := ProvidersMap[voucher.ProviderID]
			if !ok {
				continue
			}
			if voucher.TraderID == 4 {
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
				voucher.ProviderID == 24 ||
				voucher.ProviderID == 45 ||
				voucher.ProviderID == 34 {
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
	for _, bitcouTestCountry := range countriesDev {
		err = model.AddTestCountry(bitcouTestCountry)
		if err != nil {
			panic("unable to store test country information")
		}
	}
}
