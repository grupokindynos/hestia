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

	// Voucher Info by Countries
	doc := firestore.Collection("bitcou")
	docTest := firestore.Collection("bitcou_test")
	model := models.BitcouModel{Firestore: doc, FirestoreTest: docTest}

	// Bitcou Filtering System
	docFilter := firestore.Collection("bitcou_filters")
	modelFilter := models.BitcouModel{Firestore: docFilter, FirestoreTest: docFilter}

	prodProvFilter, prodVouchersFilter, err := modelFilter.GetFilters("prod")
	devProvFilter, devVoucherFilter, err := modelFilter.GetFilters("dev")

	fmt.Println(prodProvFilter, prodVouchersFilter)
	fmt.Println(devProvFilter, devVoucherFilter)

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
	fmt.Println("Missing id vouchers")
	for _, availableCountry := range availableCountry {
		newCountryData := models.BitcouCountry{
			ID:       availableCountry,
			Vouchers: []bitcou.LightVoucher{},
		}
		for _, voucher := range voucherListDev {
			_, okProv := devProvFilter[voucher.ProviderID]
			_, okVoucher := devVoucherFilter[voucher.SKU]
			available := voucher.Countries[availableCountry]
			if available && !okProv && !okVoucher {
				_, ok := devProvidersMap[voucher.ProviderID]
				if !ok {
					//fmt.Println("missing provider for: ", voucher.SKU)
					continue
				}
				newCountryData.Vouchers = append(newCountryData.Vouchers, *bitcou.NewLightVoucher(voucher))
			} else {
				log.Println("succesfully filtered ", voucher.SKU)
			}
		}
		countriesDev = append(countriesDev, newCountryData)
	}

	for _, availableCountry := range availableCountry {
		newCountryData := models.BitcouCountry{
			ID:       availableCountry,
			Vouchers: []bitcou.LightVoucher{},
		}
		for _, voucher := range voucherListProd {
			_, okProv := prodProvFilter[voucher.ProviderID]
			_, okVoucher := prodVouchersFilter[voucher.SKU]
			available := voucher.Countries[availableCountry]
			if available && !okProv && !okVoucher {
				_, ok := ProvidersMap[voucher.ProviderID]
				if !ok {
					continue
				}
				newCountryData.Vouchers = append(newCountryData.Vouchers, *bitcou.NewLightVoucher(voucher))
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
