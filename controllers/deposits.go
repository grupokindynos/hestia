package controllers

import (
	"github.com/grupokindynos/hestia/models"
)

/*

	DepositsController is a safe-access query for deposits on Firestore Database
	Database Structure:

	deposits/
		UID/
          	deposit -> Array of DepositIDs

	depositIndex/
		DepositID/
			depositData

*/

type Deposit struct{}

type DepositsController struct {
	Model *models.DepositsModel
}
