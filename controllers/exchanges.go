package controllers

import (
	"encoding/json"
	nError "errors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/common/tokens/mrt"
	"github.com/grupokindynos/common/tokens/mvt"
	"github.com/grupokindynos/hestia/models"
	"log"
	"os"
)

type ExchangesController struct {
	Model *models.ExchangesModel
}

func (ec *ExchangesController) GetExchange(c *gin.Context) {
	id := c.Query("Id")
	if id == "" {
		responses.GlobalResponseError(nil, nError.New("Missing exchange id"), c)
		return
	}
	log.Println(id)
	_, _, err := mvt.VerifyRequest(c)
	if err != nil {
		log.Println("2 " + err.Error())
		responses.GlobalResponseNoAuth(c)
		return
	}
	exchange, err := ec.Model.Get(id)
	if err != nil {
		log.Println("3 " + err.Error())
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), exchange, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	log.Println("Sale")
	return
}

func (ec *ExchangesController) GetExchanges(c *gin.Context) {
	_, _, err := mvt.VerifyRequest(c)
	if err != nil {
		log.Println(err)
		responses.GlobalResponseNoAuth(c)
		return
	}
	orders, err := ec.Model.GetAll()
	if err != nil {
		log.Println(err)
		responses.GlobalResponseError(nil, err, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), orders, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}

func (ec *ExchangesController) UpdateExchange(c *gin.Context) {
	payload, _, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}
	// Try to unmarshal the information of the payload
	var updateData hestia.ExchangeInfo
	err = json.Unmarshal(payload, &updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorUnmarshal, c)
		return
	}
	_, err = ec.Model.Get(updateData.Id)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	err = ec.Model.Update(updateData)
	if err != nil {
		responses.GlobalResponseError(nil, errors.ErrorDBStore, c)
		return
	}
	header, body, err := mrt.CreateMRTToken("hestia", os.Getenv("MASTER_PASSWORD"), updateData.Id, os.Getenv("HESTIA_PRIVATE_KEY"))
	responses.GlobalResponseMRT(header, body, c)
	return
}
