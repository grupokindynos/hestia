package controllers

import "go.mongodb.org/mongo-driver/mongo"

/*

	UsersController is a safe-access query for user information
	Database Structure:

	users/
		UID/
          	userData

*/

type User struct{}

type UsersController struct {
	DB *mongo.Database
}
