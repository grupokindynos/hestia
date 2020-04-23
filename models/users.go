package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type UsersModel struct {
	Firestore  *firestore.DocumentRef
	Collection string
}

// Get will return the user information stored on MongoDB
func (m *UsersModel) Get(uid string) (user hestia.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection).Doc(uid)
	doc, err := ref.Get(ctx)
	if err != nil {
		return user, err
	}
	err = doc.DataTo(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Update will update the user information on the MongoDB
func (m *UsersModel) Update(user hestia.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(user.ID).Set(ctx, user)
	return err
}

func (m *UsersModel) GetAll() (users []hestia.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := m.Firestore.Collection(m.Collection)
	docIterator := ref.Documents(ctx)
	docSnap, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docSnap {
		var user hestia.User
		_ = doc.DataTo(&user)
		users = append(users, user)
	}
	return users, nil
}

// AddShift will add a shift id into the user shifts array.
func (m *UsersModel) AddShift(uid string, shiftID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(uid).Update(ctx, []firestore.Update{{Path: "shifts", Value: firestore.ArrayUnion(shiftID)}})
	return err
}

// AddShift will add a shift id into the user shifts array.
func (m *UsersModel) AddShiftV2(uid string, shiftID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(uid).Update(ctx, []firestore.Update{{Path: "shift2", Value: firestore.ArrayUnion(shiftID)}})
	return err
}

// AddCard will add a card code into the user cards array.
func (m *UsersModel) AddCard(uid string, cardCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(uid).Update(ctx, []firestore.Update{{Path: "cards", Value: firestore.ArrayUnion(cardCode)}})
	return err
}

// AddVoucher will add a voucher id into the user vouchers array.
func (m *UsersModel) AddVoucher(uid string, voucherID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(uid).Update(ctx, []firestore.Update{{Path: "vouchers", Value: firestore.ArrayUnion(voucherID)}})
	return err
}

// AddDeposit will add a deposit id into the user deposits array.
func (m *UsersModel) AddDeposit(uid string, depositID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(uid).Update(ctx, []firestore.Update{{Path: "deposits", Value: firestore.ArrayUnion(depositID)}})
	return err
}

// AddOrder will add a order id into the user orders array.
func (m *UsersModel) AddOrder(uid string, orderID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Firestore.Collection(m.Collection).Doc(uid).Update(ctx, []firestore.Update{{Path: "orders", Value: firestore.ArrayUnion(orderID)}})
	return err
}
