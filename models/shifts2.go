package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"strconv"
	"time"
)

type ShiftModelV2 struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *ShiftModelV2) Get(id string) (shift hestia.ShiftV2, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return shift, err
	}
	err = doc.DataTo(&shift)
	if err != nil {
		return shift, err
	}
	return shift, nil
}

func (m *ShiftModelV2) Update(shift hestia.ShiftV2) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(shift.ID).Set(ctx, shift)
	return err
}

func (m *ShiftModelV2) GetAll(filter string, timeFilter string) (shifts []hestia.ShiftV2, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	var docSnap []*firestore.DocumentSnapshot
	if timeFilter != "" {
		timeInt, err := strconv.Atoi(timeFilter)
		if err != nil {
			return nil, err
		}
		query := ref.Where("timestamp", ">=", timeInt)
		docSnap, err = query.Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		if filter == "all" {
			docSnap, err = ref.Documents(ctx).GetAll()
			if err != nil {
				return nil, err
			}
		} else {
			query := ref.Where("status", "==", filter)
			docSnap, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, err
			}
		}
	}
	for _, doc := range docSnap {
		var shift hestia.ShiftV2
		_ = doc.DataTo(&shift)
		shifts = append(shifts, shift)
	}
	return shifts, nil
}
