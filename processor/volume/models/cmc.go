package cmc

type CoinsList struct {
	Data   []CoinInfo `json:"data"`
	Status struct {
		CreditCount  int         `json:"credit_count"`
		Elapsed      int         `json:"elapsed"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Notice       interface{} `json:"notice"`
		Timestamp    string      `json:"timestamp"`
	} `json:"status"`
}

type CoinInfo struct {
	CirculatingSupply *float64 `json:"circulating_supply"`
	CmcRank           int      `json:"cmc_rank"`
	DateAdded         string   `json:"date_added"`
	ID                int      `json:"id"`
	LastUpdated       string   `json:"last_updated"`
	MaxSupply         *float64 `json:"max_supply"`
	Name              string   `json:"name"`
	NumMarketPairs    int      `json:"num_market_pairs"`
	Platform          *struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Slug         string `json:"slug"`
		Symbol       string `json:"symbol"`
		TokenAddress string `json:"token_address"`
	} `json:"platform"`
	Quote struct {
		USD struct {
			LastUpdated      string   `json:"last_updated"`
			MarketCap        *float64 `json:"market_cap"`
			PercentChange1h  float64  `json:"percent_change_1h"`
			PercentChange24h *float64 `json:"percent_change_24h"`
			PercentChange7d  *float64 `json:"percent_change_7d"`
			Price            float64  `json:"price"`
			Volume24h        *float64 `json:"volume_24h"`
		} `json:"USD"`
	} `json:"quote"`
	Slug        string   `json:"slug"`
	Symbol      string   `json:"symbol"`
	Tags        []string `json:"tags"`
	TotalSupply *float64 `json:"total_supply"`
}
