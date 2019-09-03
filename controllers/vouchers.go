package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

/*

	VouchersController is a safe-access query for vouchers on Firestore Database
	Database Structure:

	vouchers/
		UID/
          	vouchers -> Array of VouchersIDs

	voucherIndex/
		VoucherID/
			voucherData

*/

type Voucher struct{}

type VouchersController struct {
	DB *mongo.Database
}
