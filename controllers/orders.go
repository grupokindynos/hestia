package controllers

import "cloud.google.com/go/firestore"

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
	DB *firestore.Client
}
