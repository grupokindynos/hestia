package controllers

import "cloud.google.com/go/firestore"

/*

	UsersController is a safe-access query for user information
	Database Structure:

	users/
		UID/
          	userData

*/

type User struct{}

type UsersController struct {
	DB *firestore.Client
}
