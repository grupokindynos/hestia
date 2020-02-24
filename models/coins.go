package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/grupokindynos/common/hestia"
)

type CoinsModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *CoinsModel) GetCoinsData() ([]hestia.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	var CoinData []hestia.Coin
	for _, doc := range docSnap {
		var coin hestia.Coin
		_ = doc.DataTo(&coin)
		CoinData = append(CoinData, coin)
	}
	return CoinData, nil
}

func (m *CoinsModel) UpdateCoinsData(Coins []hestia.Coin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	for _, coin := range Coins {
		docref := ref.Doc(coin.Ticker)
		_, _ = docref.Set(ctx, coin)
	}
	return nil
}

func (m *CoinsModel) AddCoin(Coin hestia.Coin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docref := ref.Doc(Coin.Ticker)
	_, _ = docref.Set(ctx, Coin)
	return nil
}
