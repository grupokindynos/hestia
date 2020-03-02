package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/hestia/services/bitcou"
	"time"
)

type BitcouCountry struct {
	ID       string           `firestore:"id" json:"id"`
	Vouchers []bitcou.Voucher `firestore:"vouchers" json:"vouchers"`
}

type BitcouFilter struct {
	ID       string           `firestore:"id" json:"id"`
	Providers []int `firestore:"providers" json:"providers"`
	Vouchers []string `firestore:"vouchers" json:"vouchers"`
}

type ApiBitcouFilter struct {
	Target string `json:"api"`
	Vouchers []string `json:"vouchers"`
	Providers []int `json:"providers"`
}

type BitcouModel struct {
	Firestore     *firestore.CollectionRef
	FirestoreTest *firestore.CollectionRef
}

type BitcouConfModel struct {
	Firestore *firestore.CollectionRef
}

func (bcm *BitcouConfModel) UpdateFilters(filter BitcouFilter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := bcm.Firestore.Doc(filter.ID).Set(ctx, filter)
	return err
}

func (bm *BitcouModel) AddTestCountry(country BitcouCountry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := bm.FirestoreTest.Doc(country.ID).Set(ctx, country)
	return err
}

func (bm *BitcouModel) AddCountry(country BitcouCountry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := bm.Firestore.Doc(country.ID).Set(ctx, country)
	return err
}

func (bm *BitcouModel) GetCountry(id string) (country BitcouCountry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := bm.Firestore.Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return country, err
	}
	err = doc.DataTo(&country)
	if err != nil {
		return country, err
	}
	return country, nil
}

func (bm *BitcouModel) GetTestCountry(id string) (country BitcouCountry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := bm.FirestoreTest.Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return country, err
	}
	err = doc.DataTo(&country)
	if err != nil {
		return country, err
	}
	return country, nil
}

func (bm *BitcouModel) GetFilters(db string) (filterMapProviders map[int]bool, filterMapVouchers  map[string]bool, err error) {
	var filter BitcouFilter
	filterMapProviders = make(map[int]bool)
	filterMapVouchers = make(map[string]bool)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := bm.Firestore.Doc(db)
	doc, err := ref.Get(ctx)
	if err != nil {
		return filterMapProviders, filterMapVouchers, err
	}
	err = doc.DataTo(&filter)
	if err != nil {
		return filterMapProviders, filterMapVouchers, err
	}

	for _, provider := range filter.Providers {
		filterMapProviders[provider] = false
	}

	for _, voucher := range filter.Vouchers {
		filterMapVouchers[voucher] = false
	}

	return filterMapProviders, filterMapVouchers, err
}
