package bitcou

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Service struct {
	BitcouURL   string
	BitcouToken string
}

func (bs *Service) GetList(dev bool) ([]Voucher, error) {
	var url string
	if dev {
		url = os.Getenv("BITCOU_URL_DEV") + "voucher/availableVouchers/"
	} else {
		url = os.Getenv("BITCOU_URL_PROD") + "voucher/availableVouchers/"
	}
	token := "Bearer " + os.Getenv("BITCOU_TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{Timeout: 300 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, _ := ioutil.ReadAll(res.Body)
	var response BaseResponse
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return nil, err
	}
	var vouchersList []Voucher
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &vouchersList)
	if err != nil {
		return nil, err
	}
	return vouchersList, nil
}

func (bs *Service) GetProviders(dev bool) ([]Provider, error) {
	var url string
	if dev {
		url = os.Getenv("BITCOU_URL_DEV") + "voucher/providers"
	} else {
		url = os.Getenv("BITCOU_URL_PROD") + "voucher/providers"
	}
	token := "Bearer " + os.Getenv("BITCOU_TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, _ := ioutil.ReadAll(res.Body)
	var response BaseResponse
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return nil, err
	}
	var providerList []Provider
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &providerList)
	if err != nil {
		return nil, err
	}
	return providerList, nil
}

func InitService() *Service {
	service := &Service{
		BitcouURL:   os.Getenv("BITCOU_URL"),
		BitcouToken: os.Getenv("BITCOU_TOKEN"),
	}
	return service
}
