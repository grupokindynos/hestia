package services

import (
	"encoding/json"
	"github.com/grupokindynos/hestia/config"
	"io/ioutil"
)

type ObolSimpleResponse struct {
	Data struct {
		AUD float64 `json:"AUD"`
		BGN float64 `json:"BGN"`
		BRL float64 `json:"BRL"`
		BTC float64 `json:"BTC"`
		CAD float64 `json:"CAD"`
		CHF float64 `json:"CHF"`
		CNY float64 `json:"CNY"`
		CZK float64 `json:"CZK"`
		DKK float64 `json:"DKK"`
		GBP float64 `json:"GBP"`
		HKD float64 `json:"HKD"`
		HRK float64 `json:"HRK"`
		HUF float64 `json:"HUF"`
		IDR float64 `json:"IDR"`
		ILS float64 `json:"ILS"`
		INR float64 `json:"INR"`
		ISK float64 `json:"ISK"`
		JPY float64 `json:"JPY"`
		KRW float64 `json:"KRW"`
		MXN float64 `json:"MXN"`
		MYR float64 `json:"MYR"`
		NOK float64 `json:"NOK"`
		NZD float64 `json:"NZD"`
		PHP float64 `json:"PHP"`
		PLN float64 `json:"PLN"`
		RON float64 `json:"RON"`
		RUB float64 `json:"RUB"`
		SEK float64 `json:"SEK"`
		SGD float64 `json:"SGD"`
		THB float64 `json:"THB"`
		TRY float64 `json:"TRY"`
		USD float64 `json:"USD"`
		ZAR float64 `json:"ZAR"`
	} `json:"data"`
	Status int `json:"status"`
}

type ObolComplexResponse struct {
	Data   float64 `json:"data"`
	Status int     `json:"status"`
}

type ObolService struct {
	URL string
}

func (o *ObolService) GetSimpleRate(coin string) (float64, error) {
	res, err := config.HttpClient.Get(o.URL + "/simple/" + coin)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var rates ObolSimpleResponse
	err = json.Unmarshal(contents, &rates)
	if err != nil {
		return 0, err
	}
	return rates.Data.BTC, nil
}

func (o *ObolService) GetComplexRate(fromcoin string, tocoin string, amount string) (float64, error) {
	res, err := config.HttpClient.Get(o.URL + "/complex/" + fromcoin + "/" + tocoin + "/" + "?amount=" + amount)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var rates ObolComplexResponse
	err = json.Unmarshal(contents, &rates)
	if err != nil {
		return 0, err
	}
	return rates.Data, nil
}
