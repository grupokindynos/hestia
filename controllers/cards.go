package controllers

import "github.com/grupokindynos/hestia/models"

/*

	CardsController is a safe-access query for cards on Firestore Database
	Database Structure:

	cards/
		UID/
          	cards -> Array of CardCodes

	cardsIndex/
		CardCode/
			carddata

	pins/
		CardCode/
			pinhash

*/

type Card struct{}

type CardsController struct {
	Model *models.CardsModel
}
