package models

import (
	"context"
	"github.com/grupokindynos/common/hestia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UsersModel struct {
	Db         *mongo.Database
	Collection string
}

// Get will return the user information stored on MongoDB
func (m *UsersModel) Get(uid string) (user hestia.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": uid}
	err = collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

// Update will update the user information on the MongoDB
func (m *UsersModel) Update(user hestia.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": user.ID}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$set", Value: user}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func (m *UsersModel) GetAll() (users []hestia.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := m.Db.Collection(m.Collection)
	curr, _ := col.Find(ctx, bson.M{})
	for curr.Next(ctx) {
		var user hestia.User
		_ = curr.Decode(&user)
		users = append(users, user)
	}
	return users, nil
}

// AddShift will add a shift id into the user shifts array.
func (m *UsersModel) AddShift(uid string, shiftID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": uid}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$push", Value: bson.M{"shifts": shiftID}}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

// AddCard will add a card code into the user cards array.
func (m *UsersModel) AddCard(uid string, cardCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": uid}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$push", Value: bson.M{"cards": cardCode}}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

// AddVoucher will add a voucher id into the user vouchers array.
func (m *UsersModel) AddVoucher(uid string, voucherID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": uid}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$push", Value: bson.M{"vouchers": voucherID}}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

// AddDeposit will add a deposit id into the user deposits array.
func (m *UsersModel) AddDeposit(uid string, depositID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": uid}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$push", Value: bson.M{"deposits": depositID}}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}

// AddOrder will add a order id into the user orders array.
func (m *UsersModel) AddOrder(uid string, orderID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": uid}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$push", Value: bson.M{"orders": orderID}}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}
