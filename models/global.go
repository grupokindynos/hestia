package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/grupokindynos/common/errors"
	"github.com/grupokindynos/common/hestia"
)

type GlobalConfigModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *GlobalConfigModel) GetConfigData() (hestia.Config, error) {
	shiftAvailable, err := m.getAvailable("shifts")
	if err != nil {
		return hestia.Config{}, errors.ErrorConfigDataGet
	}
	depositAvailable, err := m.getAvailable("deposits")
	if err != nil {
		return hestia.Config{}, errors.ErrorConfigDataGet
	}
	voucherAvailable, err := m.getAvailable("vouchers")
	if err != nil {
		return hestia.Config{}, errors.ErrorConfigDataGet
	}
	ordersAvailable, err := m.getAvailable("orders")
	if err != nil {
		return hestia.Config{}, errors.ErrorConfigDataGet
	}
	adrestiaAvailable, err := m.getAvailable("adrestia")
	if err != nil {
		return hestia.Config{}, errors.ErrorConfigDataGet
	}
	customParams, err := m.getCustomParams()
	configData := hestia.Config{
		Shift:    shiftAvailable,
		Deposits: depositAvailable,
		Vouchers: voucherAvailable,
		Orders:   ordersAvailable,
		Adrestia: adrestiaAvailable,
		Params: customParams,
	}
	return configData, nil
}

func (m *GlobalConfigModel) getAvailable(id string) (available hestia.Available, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return available, err
	}
	err = doc.DataTo(&available)
	if err != nil {
		return available, err
	}
	return available, nil
}

func (m *GlobalConfigModel) getCustomParams() (params hestia.CustomParams, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc("params")
	doc, err := ref.Get(ctx)
	if err != nil {
		return params, err
	}
	err = doc.DataTo(&params)
	if err != nil {
		return params, err
	}
	return params, nil
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
	err = m.storePropData("adrestia", config.Adrestia)
	if err != nil {
		return err
	}
	return nil
}

func (m *GlobalConfigModel) storePropData(id string, available hestia.Available) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(id).Set(ctx, available)
	return err
}
