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

	// Voucher Info by Countries
	doc := firestore.Collection("bitcou")
	docTest := firestore.Collection("bitcou_test")
	model := models.BitcouModel{Firestore: doc, FirestoreTest: docTest}

	// Bitcou Filtering System
	docFilter := firestore.Collection("bitcou_filters")
	modelFilter := models.BitcouModel{Firestore: docFilter, FirestoreTest: docFilter}

	prodFilters, err := modelFilter.GetFilters("prod")
	devFilters, err := modelFilter.GetFilters("dev")

	service := bitcou.InitService()

	prodProv, _ := service.GetProviders(false) // Retrieves public API vouchers
	var ProvidersMap = providersToMap(prodProv)

	devProv, _ := service.GetProviders(false) // Retrieves dev API vouchers
	var devProvidersMap = providersToMap(devProv)

	voucherListProd, err := service.GetList(false)
	if err != nil {
		panic("unable to load bitcou voucher list: " + err.Error())
	}
	voucherListDev, err := service.GetList(false)
	if err != nil {
		panic("unable to load bitcou voucher list")
	}
	var countries []models.BitcouCountry
	var countriesDev []models.BitcouCountry
	var availableCountry []string
	for key := range voucherListProd[0].Countries {
		availableCountry = append(availableCountry, key)
	}

	countriesDev = filterVouchersByCountry(availableCountry, voucherListDev, devFilters.ProviderFilter, devFilters.VoucherFilter, devProvidersMap)
	countries = filterVouchersByCountry(availableCountry, voucherListProd, prodFilters.ProviderFilter, devFilters.VoucherFilter, ProvidersMap)

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

func providersToMap(providers []bitcou.Provider) (providerMap map[int]string) {
	providerMap = make(map[int]string)
	for _, provider := range providers {
		providerMap[provider.Id] = provider.Name
	}
	return
}

func filterVouchersByCountry(availableCountries []string, voucherList []bitcou.Voucher, providerFilter map[int]bool, voucherFilter map[string]bool, providerMap map[int]string) []models.BitcouCountry {
	var countryInfo []models.BitcouCountry
	for _, availableCountry := range availableCountries {
		newCountryData := models.BitcouCountry{
			ID:       availableCountry,
			Vouchers: []bitcou.LightVoucher{},
		}
		for _, voucher := range voucherList {
			_, okProv := providerFilter[voucher.ProviderID]
			_, okVoucher := voucherFilter[voucher.SKU]
			available := voucher.Countries[availableCountry]
			if available && !okProv && !okVoucher {
				_, ok := providerMap[voucher.ProviderID]
				if !ok {
					//fmt.Println("missing provider for: ", voucher.SKU)
					continue
				}

				newCountryData.Vouchers = append(newCountryData.Vouchers, *bitcou.NewLightVoucher(voucher))
			} else {
				//log.Println("succesfully filtered ", voucher.SKU)
			}
		}
		countryInfo = append(countryInfo, newCountryData)
	}
	return countryInfo
}
