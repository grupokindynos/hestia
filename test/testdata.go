package test

import (
	"github.com/grupokindynos/hestia/models"
)

var TestCard = models.Card{
	Address:    "TEST-ADDRRESS",
	CardCode:   "TEST-CARDCODE",
	CardNumber: "123456778",
	City:       "TEST-CITY",
	Email:      "TEST@TEST.COM",
	FirstName:  "TEST",
	LastName:   "CARD",
	UID:        "XYZ12345678910",
}

var TestCoinData = []models.Coin{
	{Ticker: "BTC", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "LTC", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "DASH", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "POLIS", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "GRS", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "DGB", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "COLX", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "ONION", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "XSG", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "XZC", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
	{Ticker: "MNP", ShiftAvailable: false, DepositAvailable: false, VouchersAvailable: false, OrdersAvailable: false, Balances: models.Balances{HotWallet: 1, Exchanges: 1}},
}

var TestOrder = models.Order{
	ID:     "TEST-ORDER",
	UID:    "XYZ12345678910",
	Status: "COMPLETED",
	PaymentInfo: models.Payment{
		Address: "FAKE-ADDRRESS",
		Amount:  "100",
		Coin:    "POLIS",
		Txid:    "FAKE-TXID",

		Confirmations: "0",
	},
	AddressInfo: models.AddressInformation{
		City:       "FAKE-CITY",
		Country:    "FAKE-COUNTRY",
		PostalCode: "00000",
		Email:      "TEST@TEST.COM",
		Street:     "FAKE-STREET",
	},
	Delivery: models.DeliveryOption{
		Type:    "DHL",
		Service: "EXPRESS",
		DeliveryAddress: models.AddressInformation{
			City:       "FAKE-CITY",
			Country:    "FAKE-COUNTRY",
			PostalCode: "00000",
			Email:      "TEST@TEST.COM",
			Street:     "FAKE-STREET",
		},
	},
	PersonalizationData: models.PersonalizationInformation{
		PersonalData: models.PersonalInformation{
			BirthDate:   "20190211",
			CivilState:  "MARRIED",
			FirstName:   "TEST-USER",
			LastName:    "TEST-USER",
			Sex:         "MALE",
			HomeNumber:  "00-00-00-00",
			Nationality: "MEXICAN",
			PassportID:  "1234567890",
		},
		AddressData: models.AddressInformation{
			City:       "FAKE-CITY",
			Country:    "FAKE-COUNTRY",
			PostalCode: "00000",
			Email:      "TEST@TEST.COM",
			Street:     "FAKE-STREET",
		},
	},
}

var TestShift = models.Shift{
	ID:        "TEST-SHIFT",
	Status:    "COMPLETED",
	Timestamp: "000000000000",
	UID:       "XYZ12345678910",
	Payment: models.Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	Conversion: models.Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
}

var TestUser = models.User{
	ID:       "XYZ12345678910",
	Email:    "TEST@TEST.COM",
	KYCData:  models.KYCInformation{},
	Role:     "admin",
	Shifts:   []string{},
	Vouchers: []string{},
	Deposits: []string{},
	Cards:    []string{},
	Orders:   []string{},
}

var TestVoucher = models.Voucher{
	ID:         "TEST-VOUCHER",
	UID:        "XYZ12345678910",
	VoucherID:  "FAKE-VOUCHER",
	VariantID:  "FAKE-VARIANT",
	FiatAmount: "100",
	Name:       "TEST-VOUCHER",
	PaymentData: models.Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	BitcouPaymentData: models.Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	RedeemCode: "FAKE-REDEEM",
	Status:     "COMPLETED",
	Timestamp:  "00000000000",
}

var TestDeposit = models.Deposit{
	ID:  "TEST-DEPOSIT",
	UID: "XYZ12345678910",
	Payment: models.Payment{
		Address:       "FAKE-ADDR",
		Amount:        "123123123",
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: "0",
	},
	AmountInPeso: "100",
	CardCode:     "TEST-CARDCODE",
	Status:       "COMPLETED",
	Timestamp:    "000000000000",
}
