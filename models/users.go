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
	KYCData  KYCInformation `bson:"kyc_data" json:"kyd_data"`
	Role     string         `bson:"role" json:"role"`
	Shifts   []string       `bson:"shifts" json:"shifts"`
	Vouchers []string       `bson:"vouchers" json:"vouchers"`
	Deposits []string       `bson:"deposits" json:"deposits"`
	Cards    []string       `bson:"cards" json:"cards"`
}

type KYCInformation struct{}

type UsersModel struct {
	Db *mongo.Database
}

func (m *UsersModel) GetUserInformation(uid string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.Db.Collection("users")
	filter := bson.M{"_id": uid}
	err = collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (m *UsersModel) UpdateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shiftsColl := m.Db.Collection("users")
	uidFilter := bson.M{"_id": user.ID}
	upsert := true
	_, err := shiftsColl.UpdateOne(ctx, uidFilter, bson.D{{Key: "$set", Value: user}}, &options.UpdateOptions{Upsert: &upsert})
	return err
}
