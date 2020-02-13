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
	doc := firestore.Collection("bitcou")
	docTest := firestore.Collection("bitcou_test")
	model := models.BitcouModel{Firestore: doc, FirestoreTest: docTest}
	service := bitcou.InitService()

	prodProv, _ := service.GetProviders(false) // Retrieves public API vouchers
	var ProvidersMap = providersToMap(prodProv)

	devProv, _ := service.GetProviders(true) // Retrieves dev API vouchers
	var devProvidersMap = providersToMap(devProv)

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
	for key := range voucherListProd[0].Countries {
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
		fmt.Println(availableCountry)
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
	/*for _, bitcouCountry := range countries {
		err = model.AddCountry(bitcouCountry)
		if err != nil {
			panic("unable to store country information")
		}
	}*/
	for _, bitcouTestCountry := range countriesDev {
		err = model.AddTestCountry(bitcouTestCountry)
		if err != nil {
			panic("unable to store test country information")
		}
	}
}

func providersToMap(providers []bitcou.Provider) (providerMap map[int]string) {
	providerMap = make(map[int]string)
	for _, provider := range providers {
		providerMap[provider.Id] = provider.Name
	}
	return
}
