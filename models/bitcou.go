package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/hestia/services/bitcou"
	"google.golang.org/api/iterator"
	"time"
)

type BitcouCountry struct {
	ID       string                `firestore:"id" json:"id"`
	Vouchers []bitcou.LightVoucher `firestore:"vouchers" json:"vouchers"`
}

type BitcouCountryV2 struct {
	ID       string                  `firestore:"id" json:"id"`
	Vouchers []bitcou.LightVoucherV2 `firestore:"vouchers" json:"vouchers"`
}

type BitcouFilter struct {
	ID        string   `firestore:"id" json:"id"`
	Providers []int    `firestore:"providers" json:"providers"`
	Vouchers  []string `firestore:"vouchers" json:"vouchers"`
}

type ApiBitcouFilter struct {
	Target    string   `json:"api"`
	Vouchers  []string `json:"vouchers"`
	Providers []int    `json:"providers"`
}

type BitcouModel struct {
	Firestore       *firestore.CollectionRef
	FirestoreTest   *firestore.CollectionRef
	FirestoreV2     *firestore.CollectionRef
	FirestoreTestV2 *firestore.CollectionRef
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

func (bm *BitcouModel) AddTestCountryV2(country BitcouCountryV2) error {
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

func (bm *BitcouModel) AddCountryV2(country BitcouCountryV2) error {
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

func (bm *BitcouModel) GetCountryV2(id string) (country BitcouCountryV2, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := bm.FirestoreV2.Doc(id)
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

// Replaces both implementations of GetCountry
func (bm *BitcouModel) GetCountries(dev bool) (countries []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var iter *firestore.DocumentIterator
	if dev {
		iter = bm.FirestoreTest.Documents(ctx)
	} else {
		iter = bm.Firestore.Documents(ctx)
	}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return countries, err
		}
		countries = append(countries, doc.Ref.ID)
	}
	return countries, nil
}

// Replaces both implementations of GetCountry
func (bm *BitcouModel) GetCountriesV2(dev bool) (countries []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var iter *firestore.DocumentIterator
	if dev {
		iter = bm.FirestoreTestV2.Documents(ctx)
	} else {
		iter = bm.FirestoreV2.Documents(ctx)
	}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return countries, err
		}
		countries = append(countries, doc.Ref.ID)
	}
	return countries, nil
}

type BitcouFilterWrapper struct {
	ProviderFilter map[int]bool
	VoucherFilter  map[string]bool
}

func (bm *BitcouModel) GetFilters(db string) (filterResponse BitcouFilterWrapper, err error) {
	var filter BitcouFilter
	filterMapProviders := make(map[int]bool)
	filterMapVouchers := make(map[string]bool)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := bm.Firestore.Doc(db)
	doc, err := ref.Get(ctx)
	if err != nil {
		return filterResponse, err
	}
	err = doc.DataTo(&filter)
	if err != nil {
		return filterResponse, err
	}

	for _, provider := range filter.Providers {
		filterMapProviders[provider] = false
	}

	for _, voucher := range filter.Vouchers {
		filterMapVouchers[voucher] = false
	}
	filterResponse.ProviderFilter = filterMapProviders
	filterResponse.VoucherFilter = filterMapVouchers
	return filterResponse, err
}
