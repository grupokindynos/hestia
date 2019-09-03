package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Order struct{}

type OrdersModel struct {
	Db *mongo.Database
}

func (m *OrdersModel) GetOrder(depositid string) (order Order, err error) {
	return order, err
}

func (m *OrdersModel) StoreOrder(uid string, order Order) error {
	return nil
}
