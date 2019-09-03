package controllers

import (
	"github.com/grupokindynos/hestia/models"
	"go.mongodb.org/mongo-driver/mongo"
)

/*

	OrdersController is a safe-access query for orders on Firestore Database
	Database Structure:

	orders/
		UID/
          	orders -> Array of OrderIDs

	orderIndex/
		OrderID/
			orderData

*/

type Order struct{}

type OrdersController struct {
	Model *models.OrdersModel
}
