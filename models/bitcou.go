package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/hestia/services/bitcou"
	"time"
)

var CountryNames = map[string]string{}

type BitcouCountry struct {
	ID       string           `firestore:"id" json:"id"`
	Name     string           `firestore:"name" json:"name"`
	Vouchers []bitcou.Voucher `firestore:"vouchers" json:"vouchers"`
}

type BitcouModel struct {
	Firestore *firestore.CollectionRef
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
