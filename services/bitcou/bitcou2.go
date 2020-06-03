package bitcou

import "os"

type ServiceV2 struct {
	BitcouURL   string
	BitcouToken string
}

func InitServiceV2(url string) *ServiceV2 {
	service := &ServiceV2{
		BitcouURL:   url,
		BitcouToken: os.Getenv("BITCOU_TOKEN_V2"),
	}
	return service
}