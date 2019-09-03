package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct{}

type UsersModel struct {
	Db *mongo.Database
}

func (m *UsersModel) GetVoucher(uid string) (user User, err error) {
	return user, err
}

func (m *UsersModel) UpdateUser(user User) error {
	return nil
}
