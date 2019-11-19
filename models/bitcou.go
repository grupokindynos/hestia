package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"time"
)

type BitcouCountries struct {
	Countries []string `firestore:"countries" json:"countries"`
}

type BitcouModel struct {
	Firestore *firestore.DocumentRef
}

func (bm *BitcouModel) StoreCountries(countries BitcouCountries) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := bm.Firestore.Set(ctx, countries)
	return err
}

