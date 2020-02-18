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
	Countries    map[string]bool `firestore:"countries" json:"countries"`
	Image        string          `firestore:"image" json:"image"`
	Name         string          `firestore:"name" json:"name"`
	ProductID    int             `firestore:"product_id" json:"product_id"`
	RedeemPlace  RedeemPlace     `firestore:"redeem_place" json:"redeem_place"`
	Shipping     Shipping        `firestore:"shipping" json:"shipping"`
	TraderID     int             `firestore:"trader_id" json:"trader_id"`
	Variants     []Variants      `firestore:"variants" json:"variants"`
	ProviderID   int             `firestore:"provider_id" json:"provider_id"`
	ProviderName string          `firestore:"provider_name" json:"provider_name"`
	Benefits     map[string]bool `firestore:"benefits" json:"benefits"`
	Description	 string			 `firestore:"description" json:"description"`
	Valid	 int64			 `firestore:"valid" json:"valid"`
}

type BaseResponse struct {
	Data []interface{} `json:"data"`
	Meta MetaData      `json:"meta"`
}

type MetaData struct {
	Datetime string `json:"datetime"`
}

type Provider struct {
	Id int `json:"provider_id"`
	Name string `json:"provider_name"`
}
