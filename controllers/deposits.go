package controllers

import "cloud.google.com/go/firestore"

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
	DB *firestore.Client
}
