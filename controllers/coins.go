package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
	"os"
)

/*

	CoinsController is a safe-access query for cards on Firestore Database
	Database Structure:

	coins/
		TICKER/
          	coinAvailability

*/

type CoinsController struct {
	Model *models.CoinsModel
}

func (cc *CoinsController) GetCoinsAvailability(userData hestia.User, params Params) (interface{}, error) {
load:
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		return nil, errors.ErrorCoinDataGet
	}
	// First we check if CoinsData contains all coins on CoinFactory
	coinsDataMap := make(map[string]hestia.Coin)
	for _, coin := range coins {
		coinsDataMap[coin.Ticker] = coin
	}
	requireUpdate := false
	for k, v := range coinfactory.Coins {
		_, ok := coinsDataMap[k]
		if !ok {
			requireUpdate = true
			// If doesn't exists it means we must create it.
			newCoinData := hestia.Coin{
				Ticker: v.Tag,
				Shift: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Deposits: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Vouchers: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Orders: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Balances: hestia.BalanceLimits{
					HotWallet: 0,
					Exchanges: 0,
				},
			}
			coins = append(coins, newCoinData)
		}
	}
	if requireUpdate {
		err := cc.Model.UpdateCoinsData(coins)
		if err != nil {
			return nil, errors.ErrorCoinDataGet
		}
		goto load
	}
	return coins, nil
}

func (cc *CoinsController) GetCoinsAvailabilityMicroService(c *gin.Context) {
	_, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
load:
	coins, err := cc.Model.GetCoinsData()
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	// First we check if CoinsData contains all coins on CoinFactory
	coinsDataMap := make(map[string]hestia.Coin)
	for _, coin := range coins {
		coinsDataMap[coin.Ticker] = coin
	}
	requireUpdate := false
	for k, v := range coinfactory.Coins {
		_, ok := coinsDataMap[k]
		if !ok {
			requireUpdate = true
			// If doesn't exists it means we must create it.
			newCoinData := hestia.Coin{
				Ticker: v.Tag,
				Shift: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Deposits: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Vouchers: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Orders: hestia.Properties{
					FeePercentage: 10,
					Available:     false,
				},
				Balances: hestia.BalanceLimits{
					HotWallet: 0,
					Exchanges: 0,
				},
			}
			coins = append(coins, newCoinData)
		}
	}
	if requireUpdate {
		err := cc.Model.UpdateCoinsData(coins)
		if err != nil {
			responses.GlobalResponseError(nil, err, c)
		}
		goto load
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), coins, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (cc *CoinsController) UpdateCoinsAvailability(userData hestia.User, params Params) (interface{}, error) {
	var newCoinsData []hestia.Coin
	err := json.Unmarshal(params.Body, &newCoinsData)
	if err != nil {
		return nil, errors.ErrorUnmarshal
	}
	err = cc.Model.UpdateCoinsData(newCoinsData)
	if err != nil {
		return nil, errors.ErrorDBStore
	}
	return true, nil
}
