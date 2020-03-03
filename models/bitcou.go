package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
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
	Vouchers []int `firestore:"vouchers" json:"vouchers"`
}

type BitcouModel struct {
	Firestore     *firestore.CollectionRef
	FirestoreTest *firestore.CollectionRef
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
	fmt.Println(doc)
	if err != nil {
		return country, err
	}
	return country, nil
}

func (bm *BitcouModel) GetFilters(db string) (filterMapProviders map[int]bool, filterMapVouchers  map[int]bool, err error) {
	var filter BitcouFilter
	filterMapProviders = make(map[int]bool)
	filterMapVouchers = make(map[int]bool)

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
