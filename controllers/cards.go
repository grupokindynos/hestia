package controllers

import "go.mongodb.org/mongo-driver/mongo"

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
	DB *mongo.Database
}
