package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/hestia/config"
	"time"
)

type GlobalConfigModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *GlobalConfigModel) GetConfigData() (hestia.Config, error) {
	shiftsProps, err := m.getPropData("shifts")
	if err != nil {
		return hestia.Config{}, config.ErrorConfigDataGet
	}
	depositsProps, err := m.getPropData("deposits")
	if err != nil {
		return hestia.Config{}, config.ErrorConfigDataGet
	}
	vouchersProps, err := m.getPropData("vouchers")
	if err != nil {
		return hestia.Config{}, config.ErrorConfigDataGet
	}
	ordersProps, err := m.getPropData("orders")
	if err != nil {
		return hestia.Config{}, config.ErrorConfigDataGet
	}
	configData := hestia.Config{
		Shift:    shiftsProps,
		Deposits: depositsProps,
		Vouchers: vouchersProps,
		Orders:   ordersProps,
	}
	return configData, nil
}

func (m *GlobalConfigModel) getPropData(id string) (props hestia.Properties, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return props, err
	}
	err = doc.DataTo(&props)
	if err != nil {
		return props, err
	}
	return props, nil
}

func (m *GlobalConfigModel) UpdateConfigData(config hestia.Config) error {
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

func (m *GlobalConfigModel) storePropData(id string, props hestia.Properties) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection("polispay").Doc("hestia").Collection(m.Collection).Doc(id).Set(ctx, props)
	return err
}
