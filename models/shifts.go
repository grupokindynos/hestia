package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type ShiftModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

func (m *ShiftModel) Get(id string) (shift hestia.Shift, err error) {
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

func (m *ShiftModel) Update(shift hestia.Shift) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(shift.ID).Set(ctx, shift)
	return err
}

func (m *ShiftModel) GetAll(filter string) (shifts []hestia.Shift, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var shift hestia.Shift
		_ = doc.DataTo(&shift)
		if filter == "all" {
			shifts = append(shifts, shift)
		} else {
			if shift.Status == filter {
				shifts = append(shifts, shift)
			}
		}
	}
	return shifts, nil
}
