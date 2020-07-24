package models

import (
	"fmt"
	"log"
	"testing"
)

func TestAdrestiaModel_GetBalancers(t *testing.T) {
	balanceData, err := adrestiaModel.GetBalancers(false)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("%+v\n", balanceData)
}

func TestAdrestiaModel_GetWithdrawals(t *testing.T) {
	withdrawals, err := adrestiaModel.GetAllSimpleTx(false, 0, "withdrawals", "")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%+v\n", withdrawals)
}
