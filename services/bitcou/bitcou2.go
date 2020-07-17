package bitcou

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"github.com/grupokindynos/common/ladon"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ServiceV2 struct {
	BitcouURL   string
	BitcouToken string
	ImageMap map[int]ProviderImage
}

func InitServiceV2() *ServiceV2 {
	service := &ServiceV2{
		BitcouURL:   os.Getenv("BITCOU_URL_DEV_V2"),
		BitcouToken: os.Getenv("BITCOU_TOKEN_V2"),
		ImageMap: make(map[int]ProviderImage),
	}
	return service
}

func (bs *ServiceV2) GetListV2(dev bool) ([]VoucherV2, error) {
	var url string
	if dev {
		url = os.Getenv("BITCOU_URL_DEV_V2") + "voucher/availableVouchers/"
	} else {
		url = os.Getenv("BITCOU_URL_PROD_V2") + "voucher/availableVouchers/"
	}
	log.Println("Getting products using url: ", url)
	token := "Bearer " + os.Getenv("BITCOU_TOKEN_V2")
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
	var vouchersList []VoucherV2
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

func (bs *ServiceV2) GetProvidersV2(dev bool) ([]Provider, error) {
	var url string
	if dev {
		url = os.Getenv("BITCOU_URL_DEV_V2") + "voucher/providers"
	} else {
		url = os.Getenv("BITCOU_URL_PROD_V2") + "voucher/providers"
	}
	log.Println("Getting providers using url: ", url)
	token := "Bearer " + os.Getenv("BITCOU_TOKEN_V2")
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

func (bs *ServiceV2) GetProviderImage(providerId int, dev bool) (imageInfo ProviderImage, err error) {
	if val, ok := bs.ImageMap[providerId]; ok {
		//log.Println("using cached image for ", providerId)
		return val, nil
	}
	var url string
	if dev {
		url = os.Getenv("BITCOU_URL_DEV_V2") + "voucher/providerImage"
	} else {
		url = os.Getenv("BITCOU_URL_PROD_V2") + "voucher/providerImage"
	}
	token := "Bearer " + os.Getenv("BITCOU_TOKEN_V2")
	req, err := http.NewRequest("GET", url + "?provider_id=" + strconv.Itoa(providerId), nil)
	if err != nil {
		return imageInfo, err
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return imageInfo, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, _ := ioutil.ReadAll(res.Body)
	var response BaseResponse
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return imageInfo, err
	}
	if len(response.Data) > 0 {
		dataBytes, err := json.Marshal(response.Data[0])
		if err != nil {
			return imageInfo, err
		}
		err = json.Unmarshal(dataBytes, &imageInfo)
		if err != nil {
			return imageInfo, err
		}
		bs.ImageMap[providerId] = imageInfo
		return imageInfo, nil
	} else {
		return imageInfo, errors.New("image unavailable")
	}
}

func (bs *ServiceV2) GetProviderImageBase64(imageUrl string, providerId int) (imageInfo ladon.ProviderImageApp, err error) {
	token := "Bearer " + os.Getenv("BITCOU_TOKEN_V2")
	req, err := http.NewRequest("GET", imageUrl, nil)
	if err != nil {
		return imageInfo, err
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return imageInfo, err
	}
	contents, _ := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	uEnc := b64.StdEncoding.EncodeToString(contents)
	imageInfo.Image = uEnc
	imageInfo.ProviderId = providerId
	imageInfo.Url = imageUrl
	return
}