package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type User struct {
	ID       string         `bons:"id" json:"id"`
	Email    string         `bson:"email" json:"email"`
	KYCData  KYCInformation `bson:"kyc_data" json:"kyc_data"`
	Role     string         `bson:"role" json:"role"`
	Shifts   []string       `bson:"shifts" json:"shifts"`
	Vouchers []string       `bson:"vouchers" json:"vouchers"`
	Deposits []string       `bson:"deposits" json:"deposits"`
	Cards    []string       `bson:"cards" json:"cards"`
	Orders   []string       `bson:"orders" json:"orders"`
}

type KYCInformation struct{}

type UsersModel struct {
	Db         *mongo.Database
	Collection string
}

// GetUserInformation will return the user information stored on MongoDB
func (m *UsersModel) GetUserInformation(uid string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.Db.Collection(m.Collection)
	filter := bson.M{"_id": uid}
	err = collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

// UpdateUser will update the user information on the MongoDB
func (m *UsersModel) UpdateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection(m.Collection)
	uidFilter := bson.M{"_id": user.ID}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$set", Value: user}}, &options.UpdateOptions{Upsert: &upsert})
	return err
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
