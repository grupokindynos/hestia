package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
	"errors"
)

type AdrestiaModel struct {
	Firestore  *firestore.DocumentRef
	Collections map[string]string
}

func NewAdrestiaModel(doc firestore.DocumentRef) AdrestiaModel {
	adrestia := AdrestiaModel{Firestore: &doc, Collections: map[string]string{
		"withdrawals": "adrestia_withdrawals",
		"deposits": "adrestia_deposits",
		"orders": "adrestia_orders",
		"balancer": "adrestia_balancer",
	}}
	return adrestia
}

func (am *AdrestiaModel) GetSimpleTx(id string, txType string) (simpleTx hestia.SimpleTx, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := am.Firestore.Collection(am.Collections[txType]).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return simpleTx, err
	}
	err = doc.DataTo(&simpleTx)
	if err != nil {
		return simpleTx, err
	}
	return simpleTx, nil
}

func (am *AdrestiaModel) GetBalancerOrder(id string) (order hestia.BalancerOrder, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := am.Firestore.Collection(am.Collections["orders"]).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return order, err
	}
	err = doc.DataTo(&order)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (am *AdrestiaModel) GetBalancer(id string) (balancer hestia.Balancer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := am.Firestore.Collection(am.Collections["balancer"]).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return balancer, err
	}
	err = doc.DataTo(&balancer)
	if err != nil {
		return balancer, err
	}
	return balancer, nil
}

func (am *AdrestiaModel) UpdateSimpleTx(simpleTx hestia.SimpleTx, txType string) error {
	collection, ok := am.Collections[txType]
	if !ok {
		return errors.New("tx type not found " + txType)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := am.Firestore.Collection(collection).Doc(simpleTx.Id).Set(ctx, simpleTx)
	return err
}

func (am *AdrestiaModel) UpdateBalancerOrder(order hestia.BalancerOrder) error {
	collection := am.Collections["orders"]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := am.Firestore.Collection(collection).Doc(order.Id).Set(ctx, order)
	return err
}

func (am *AdrestiaModel) UpdateBalancer(balancer hestia.Balancer) error {
	collection := am.Collections["balancer"]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := am.Firestore.Collection(collection).Doc(balancer.Id).Set(ctx, balancer)
	return err
}

func (am *AdrestiaModel) GetAllSimpleTx(includeComplete bool, sinceTimestamp int, txType string) (simpleTxs []hestia.SimpleTx, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := am.Firestore.Collection(am.Collections[txType])
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var simpleTx hestia.SimpleTx
		_ = doc.DataTo(&simpleTx)
		if sinceTimestamp != 0 {
			if int(simpleTx.CreatedTime) >= sinceTimestamp {
				if includeComplete {
					simpleTxs = append(simpleTxs, simpleTx)
				} else if simpleTx.Status != hestia.SimpleTxStatusCompleted {
					simpleTxs = append(simpleTxs, simpleTx)
				}
			}
		} else {
			if includeComplete {
				simpleTxs = append(simpleTxs, simpleTx)
			} else if simpleTx.Status != hestia.SimpleTxStatusCompleted {
				simpleTxs = append(simpleTxs, simpleTx)
			}
		}
	}
	return simpleTxs, nil
}

func (am *AdrestiaModel) GetAllBalancerOrder(includeComplete bool, sinceTimestamp int) (orders []hestia.BalancerOrder, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := am.Firestore.Collection(am.Collections["orders"])
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var order hestia.BalancerOrder
		_ = doc.DataTo(&order)
		if sinceTimestamp != 0 {
			if int(order.CreatedTime) >= sinceTimestamp {
				if includeComplete {
					orders = append(orders, order)
				} else if order.Status != hestia.BalancerOrderStatusCompleted {
					orders = append(orders, order)
				}
			}
		} else {
			if includeComplete {
				orders = append(orders, order)
			} else if order.Status != hestia.BalancerOrderStatusCompleted {
				orders = append(orders, order)
			}
		}
	}
	return orders, nil
}

func (am *AdrestiaModel) GetAllBalancer(includeComplete bool, sinceTimestamp int) (balancers []hestia.Balancer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := am.Firestore.Collection(am.Collections["balancer"])
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var balancer hestia.Balancer
		_ = doc.DataTo(&balancer)
		if sinceTimestamp != 0 {
			if int(balancer.CreatedTime) >= sinceTimestamp {
				if includeComplete {
					balancers = append(balancers, balancer)
				} else if balancer.Status != hestia.BalancerStatusCompleted {
					balancers = append(balancers, balancer)
				}
			}
		} else {
			if includeComplete {
				balancers = append(balancers, balancer)
			} else if balancer.Status != hestia.BalancerStatusCompleted {
				balancers = append(balancers, balancer)
			}
		}
	}
	return balancers, nil
}