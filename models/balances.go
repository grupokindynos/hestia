package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type BalancesModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *BalancesModel) GetBalances() ([]hestia.CoinBalances, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	var BalancesData []hestia.CoinBalances
	for _, doc := range docSnap {
		var balance hestia.CoinBalances
		_ = doc.DataTo(&balance)
		BalancesData = append(BalancesData, balance)
	}
	return BalancesData, nil
}

func (m *BalancesModel) UpdateBalances(Balances []hestia.CoinBalances) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	for _, coin := range Balances {
		docref := ref.Doc(coin.Ticker)
		_, _ = docref.Set(ctx, coin)
	}
	return nil
}
