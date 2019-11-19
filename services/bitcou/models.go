package bitcou

type Variants struct {
	Currency  string  `json:"currency"`
	Ean       string  `json:"ean"`
	Price     float64 `json:"price"`
	Value     float64 `json:"value"`
	VariantID string  `json:"variant_id"`
}

type RedeemPlace struct {
	Market bool `json:"market"`
	Online bool `json:"online"`
}

type Shipping struct {
	EMail bool `json:"e_mail"`
	Mail  bool `json:"mail"`
	Print bool `json:"print"`
}

type Voucher struct {
	Countries   map[string]bool `json:"countries"`
	Image       string          `json:"image"`
	Name        string          `json:"name"`
	ProductID   int             `json:"product_id"`
	RedeemPlace RedeemPlace     `json:"redeem_place"`
	Shipping    Shipping        `json:"shipping"`
	TraderID    int             `json:"trader_id"`
	Variants    []Variants      `json:"variants"`
	Benefits    map[string]bool `json:"benefits"`
}

type BaseResponse struct {
	Data []interface{} `json:"data"`
	Meta MetaData      `json:"meta"`
}

type MetaData struct {
	Datetime string `json:"datetime"`
}
