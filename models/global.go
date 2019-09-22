package models

import (
	"context"
	"github.com/grupokindynos/hestia/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Properties struct {
	FeePercentage int  `bson:"fee_percentage" json:"fee_percentage"`
	Available     bool `bson:"available" json:"available"`
}

type Config struct {
	Shift    Properties `bson:"shift" json:"shift"`
	Deposits Properties `bson:"deposits" json:"deposits"`
	Vouchers Properties `bson:"vouchers" json:"vouchers"`
	Orders   Properties `bson:"orders" json:"orders"`
}

type GlobalConfigModel struct {
	Db         *mongo.Database
	Collection string
}

func (m *GlobalConfigModel) GetConfigData() (Config, error) {
	shiftsProps, err := m.getPropData("shifts")
	if err != nil {
		return Config{}, config.ErrorConfigDataGet
	}
	depositsProps, err := m.getPropData("deposits")
	if err != nil {
		return Config{}, config.ErrorConfigDataGet
	}
	vouchersProps, err := m.getPropData("vouchers")
	if err != nil {
		return Config{}, config.ErrorConfigDataGet
	}
	ordersProps, err := m.getPropData("orders")
	if err != nil {
		return Config{}, config.ErrorConfigDataGet
	}
	configData := Config{
		Shift:    shiftsProps,
		Deposits: depositsProps,
		Vouchers: vouchersProps,
		Orders:   ordersProps,
	}
	return configData, nil
}

func (m *GlobalConfigModel) getPropData(id string) (props Properties, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	err = col.FindOne(ctx, filter).Decode(&props)
	return props, err
}

func (m *GlobalConfigModel) UpdateConfigData(config Config) error {
	err := m.storePropData("shifts", config.Shift)
	if err != nil {
		return err
	}
	err = m.storePropData("deposits", config.Deposits)
	if err != nil {
		return err
	}
	err = m.storePropData("vouchers", config.Vouchers)
	if err != nil {
		return err
	}
	err = m.storePropData("orders", config.Orders)
	if err != nil {
		return err
	}
	return nil
}

func (m *GlobalConfigModel) storePropData(id string, props Properties) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": id}
	upsert := true
	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: props}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}
