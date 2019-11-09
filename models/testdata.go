package models

import (
	"github.com/grupokindynos/common/hestia"
	"time"
)

var TestFeePayment = hestia.Payment{
	Address:       "TEST-ADDR",
	Amount:        0,
	Coin:          "polis",
	RawTx:         "TEST-RAW-TX",
	Txid:          "TEST-TXID",
	Confirmations: 0,
}

var TestCard = hestia.Card{
	Address:    "TEST-ADDRRESS",
	CardCode:   "TEST-CARDCODE",
	CardNumber: "123456778",
	City:       "TEST-CITY",
	Email:      "TEST@TEST.COM",
	FirstName:  "TEST",
	LastName:   "CARD",
	UID:        "XYZ12345678910",
}

var TestCoinData = []hestia.Coin{
	{Ticker: "BTC", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "COLX", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "DASH", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "DGB", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "ETH", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "GRS", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "LTC", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "ONION", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "POLIS", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "TELOS", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "TUSD", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "USDC", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "USDT", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "XSG", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
	{Ticker: "XZC", Shift: hestia.Properties{FeePercentage: 10, Available: true}, Deposits: hestia.Properties{FeePercentage: 10, Available: true}, Vouchers: hestia.Properties{FeePercentage: 10, Available: true}, Orders: hestia.Properties{FeePercentage: 10, Available: true}, Balances: hestia.BalanceLimits{HotWallet: 0, Exchanges: 0}},
}

var TestBalances = []hestia.CoinBalances{
	{Ticker: "BTC", Balance: 0, Status: "SUCCESS"},
	{Ticker: "COLX", Balance: 0, Status: "SUCCESS"},
	{Ticker: "DASH", Balance: 0, Status: "SUCCESS"},
	{Ticker: "DGB", Balance: 0, Status: "SUCCESS"},
	{Ticker: "ETH", Balance: 0, Status: "SUCCESS"},
	{Ticker: "GRS", Balance: 0, Status: "SUCCESS"},
	{Ticker: "LTC", Balance: 0, Status: "SUCCESS"},
	{Ticker: "ONION", Balance: 0, Status: "SUCCESS"},
	{Ticker: "POLIS", Balance: 0, Status: "SUCCESS"},
	{Ticker: "TELOS", Balance: 0, Status: "SUCCESS"},
	{Ticker: "TUSD", Balance: 0, Status: "SUCCESS"},
	{Ticker: "USDC", Balance: 0, Status: "SUCCESS"},
	{Ticker: "USDT", Balance: 0, Status: "SUCCESS"},
	{Ticker: "XSG", Balance: 0, Status: "SUCCESS"},
	{Ticker: "XZC", Balance: 0, Status: "SUCCESS"},
}

var TestConfigData = hestia.Config{
	Shift:    hestia.Available{Available: true},
	Deposits: hestia.Available{Available: true},
	Orders:   hestia.Available{Available: true},
	Vouchers: hestia.Available{Available: true},
}

var TestOrder = hestia.Order{
	ID:     "TEST-ORDER",
	UID:    "XYZ12345678910",
	Status: "COMPLETED",
	PaymentInfo: hestia.Payment{
		Address: "FAKE-ADDRRESS",
		Amount:  100,
		Coin:    "POLIS",
		Txid:    "FAKE-TXID",

		Confirmations: 0,
	},
	FeePayment: TestFeePayment,
	AddressInfo: hestia.AddressInformation{
		City:       "FAKE-CITY",
		Country:    "FAKE-COUNTRY",
		PostalCode: "00000",
		Email:      "TEST@TEST.COM",
		Street:     "FAKE-STREET",
	},
	Delivery: hestia.DeliveryOption{
		Type:    "DHL",
		Service: "EXPRESS",
		DeliveryAddress: hestia.AddressInformation{
			City:       "FAKE-CITY",
			Country:    "FAKE-COUNTRY",
			PostalCode: "00000",
			Email:      "TEST@TEST.COM",
			Street:     "FAKE-STREET",
		},
	},
	PersonalizationData: hestia.PersonalizationInformation{
		PersonalData: hestia.PersonalInformation{
			BirthDate:   "20190211",
			CivilState:  "MARRIED",
			FirstName:   "TEST-USER",
			LastName:    "TEST-USER",
			Sex:         "MALE",
			HomeNumber:  "00-00-00-00",
			Nationality: "MEXICAN",
			PassportID:  "1234567890",
		},
		AddressData: hestia.AddressInformation{
			City:       "FAKE-CITY",
			Country:    "FAKE-COUNTRY",
			PostalCode: "00000",
			Email:      "TEST@TEST.COM",
			Street:     "FAKE-STREET",
		},
	},
}

var TestShift = hestia.Shift{
	ID:        "TEST-SHIFT",
	Status:    hestia.GetShiftStatusString(hestia.ShiftStatusComplete),
	Timestamp: "000000000000",
	UID:       "XYZ12345678910",
	Payment: hestia.Payment{
		Address:       "FAKE-ADDR",
		Amount:        0123123,
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: 0,
	},
	FeePayment: TestFeePayment,
	Rate: hestia.ShiftRate{
		Rate:       0.1,
		FromCoin:   "polis",
		FromAmount: 1000,
		ToCoin:     "dash",
		ToAmount:   1000,
		ToAddress:  "FakE_ADDR",
		FeeCoin:    "polis",
		FeeAmount:  1000,
		FeeAddress: "FAKE_FEE_ADDR",
	},
}

var TestUser = hestia.User{
	ID:       "XYZ12345678910",
	Email:    "TEST@TEST.COM",
	KYCData:  hestia.KYCInformation{},
	Role:     "admin",
	Shifts:   []string{},
	Vouchers: []string{},
	Deposits: []string{},
	Cards:    []string{},
	Orders:   []string{},
}

var TestVoucher = hestia.Voucher{
	ID:         "TEST-VOUCHER",
	UID:        "XYZ12345678910",
	VoucherID:  123123,
	VariantID:  "FAKE-VARIANT",
	FiatAmount: 123123,
	Name:       "TEST-VOUCHER",
	FeePayment: TestFeePayment,
	PaymentData: hestia.Payment{
		Address:       "FAKE-ADDR",
		Amount:        123123,
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: 123,
	},
	BitcouPaymentData: hestia.Payment{
		Address:       "FAKE-ADDR",
		Amount:        123123,
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: 123,
	},
	BitcouID:        "FAKE-ID",
	RedeemCode:      "FAKE-REDEEM",
	Status:          hestia.GetVoucherStatusString(hestia.VoucherStatusComplete),
	Timestamp:       time.Unix(0, 0).Unix(),
	RedeemTimestamp: time.Unix(0, 0).Unix(),
}

var TestDeposit = hestia.Deposit{
	ID:  "TEST-DEPOSIT",
	UID: "XYZ12345678910",
	Payment: hestia.Payment{
		Address:       "FAKE-ADDR",
		Amount:        123123,
		Coin:          "POLIS",
		Txid:          "FAKE-TXID",
		Confirmations: 123,
	},
	FeePayment:   TestFeePayment,
	AmountInPeso: "100",
	CardCode:     "TEST-CARDCODE",
	Status:       "COMPLETED",
	Timestamp:    "000000000000",
}

var TestExchangeData = hestia.AdrestiaOrder{
	ID:              "TEST-ORDER",
	Exchange:        "fake-exchange",
	Time:            0000000000,
	Status:          hestia.GetAdrestiaStatusString(hestia.AdrestiaStatusCompleted),
	Amount:          10000,
	FromCoin:        "fake-coin",
	ToCoin:          "fake-coin",
	WithdrawAddress: "FAKE-ADDR",
	Message:         "NO-MESSAGE",
}
