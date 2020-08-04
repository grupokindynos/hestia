package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/grupokindynos/common/ladon"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services/bitcou"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var (
	service = bitcou.InitServiceV2()
	providerImages = make(map[int]ladon.ProviderImageApp)
)

// This tool must be run every 12 hours to index the bitcou vouchers list.
func main() {
	_ = godotenv.Load()

	// TODO Firebase Shit
	model, prodFilter, devFilter := GetFirebaseData() // DB model, and voucher filters


	prodProv, _ := service.GetProvidersV2(false) // Retrieves public API vouchers
	var ProvidersMap = providersToMap(prodProv)

	devProv, _ := service.GetProvidersV2(false) // Retrieves dev API vouchers
	var devProvidersMap = providersToMap(devProv)

	voucherListProd, err := service.GetListV2(false)
	if err != nil {
		panic("unable to load bitcou voucher list: " + err.Error())
	}
	voucherListDev, err := service.GetListV2(true)
	if err != nil {
		panic("unable to load bitcou voucher list: " + err.Error())
	}

	file, _ := json.MarshalIndent(voucherListDev, "", " ")
	err = ioutil.WriteFile("vouchers.json", file, 0644)
	if err != nil {
		fmt.Println(err)
	}
	var countries []models.BitcouCountryV2
	var countriesDev []models.BitcouCountryV2

	// Voucher Filter
	countries = filterVouchersByCountry(voucherListProd, prodFilter.ProviderFilter, prodFilter.VoucherFilter, ProvidersMap)
	countriesDev = filterVouchersByCountry(voucherListDev, devFilter.ProviderFilter, devFilter.VoucherFilter, devProvidersMap)


	for _, bitcouCountry := range countries {
		if bitcouCountry.ID == "BR" {
			log.Println("Ignoring country ", bitcouCountry.ID)
		} else {
			err = model.AddCountryV2(bitcouCountry)
			if err != nil {
				panic("unable to store country information")
			}
		}
		for _, product := range bitcouCountry.Vouchers {
			if _, ok := providerImages[product.ProviderID]; !ok {
				imageBase64, err := service.GetProviderImageBase64(product.Image, product.ProviderID)
				if err != nil {
					continue
				}
				providerImages[product.ProviderID] = imageBase64
			}
		}
	}

	for _, imageInfo := range providerImages {
		err := model.AddProviderImage(imageInfo)
		if err != nil {
			log.Println(err)
		}
	}

	for _, bitcouTestCountry := range countriesDev { 
		err = model.AddTestCountryV2(bitcouTestCountry)
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

func filterVouchersByCountry(voucherList []bitcou.VoucherV2, providerFilter map[int]bool, voucherFilter map[string]bool, providerMap map[int]string) []models.BitcouCountryV2 {
	var countryInfo []models.BitcouCountryV2
	countryMap := make(map[string]models.BitcouCountryV2)

	for _, voucher := range voucherList {
		strId := strconv.Itoa(voucher.ProductID)
		_, okProv := providerFilter[voucher.ProviderID]
		_, okVoucher := voucherFilter[strId]
		if voucher.ProductID == 55 {
			log.Println(voucher.ProviderName, "!")
		}
		if !okProv && !okVoucher {
			_, ok := providerMap[voucher.ProviderID]
			if !ok {
				fmt.Println("missing provider for: ", voucher.ProductID)
				continue
			}
			imageInfo, err := service.GetProviderImage(voucher.ProviderID, false)
			var imageStr string
			if err != nil {
				imageStr = "unknown"
			} else {
				if imageInfo.Image == "" {
					imageStr = "unknown"
				} else {
					imageStr = imageInfo.Image
				}
			}
			for _, country := range voucher.Countries {
				countryData, ok := countryMap[country]
				if !ok {
					newCountry := models.BitcouCountryV2{
						ID:       country,
						Vouchers: []bitcou.LightVoucherV2{},
					}
					newVoucherV2 := bitcou.NewLightVoucherV2(voucher, imageStr)
					if len(newVoucherV2.Variants) > 0 {
						newCountry.Vouchers = append(newCountry.Vouchers, *newVoucherV2)
						countryMap[country] = newCountry
					} else {
						log.Println(voucher.ProviderName, " has no variants left")
					}
				} else {
					newVoucherV2 := bitcou.NewLightVoucherV2(voucher, imageStr)
					if len(newVoucherV2.Variants) > 0 {
						countryData.Vouchers = append(countryData.Vouchers, *newVoucherV2)
						countryMap[country] = countryData
					} else {
						log.Println(voucher.ProviderName, " has no variants left")
					}
				}
			}
		} else {
			//log.Println("succesfully filtered ", voucher.SKU)
		}
	}

	// TODO Map to List
	for _, countryData := range countryMap {
		countryInfo = append(countryInfo, countryData)
	}
	return countryInfo
}

func GetFirebaseData() (models.BitcouModel, models.BitcouFilterWrapper, models.BitcouFilterWrapper) {
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
	doc := firestore.Collection("bitcou2")
	docTest := firestore.Collection("bitcou_test2")
	docImages := firestore.Collection("bitcou_images")
	model := models.BitcouModel{Firestore: doc, FirestoreTest: docTest, ProductImages: docImages}

	// Bitcou Filtering System
	docFilter := firestore.Collection("bitcou_filters")
	modelFilter := models.BitcouModel{Firestore: docFilter, FirestoreTest: docFilter}

	prodFilter, err := modelFilter.GetFilters("prod")
	devFilter, err := modelFilter.GetFilters("dev")
	return model, prodFilter, devFilter
}
