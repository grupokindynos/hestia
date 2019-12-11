package models

type ShiftStatsResponse struct {
	Pending    int     `json:"pending"`
	Confirming int     `json:"confirming"`
	Confirmed  int     `json:"confirmed"`
	Error      int     `json:"error"`
	Refund     int     `json:"refund"`
	Refunded   int     `json:"refunded"`
	Complete   int     `json:"complete"`
	Total      int     `json:"total"`
	Volume     float64 `json:"volume"`
}

type VoucherStatsResponse struct {
	Pending           int     `json:"pending"`
	Confirmed         int     `json:"confirmed"`
	Confirming        int     `json:"confirming"`
	AwaitingProvider  int     `json:"awaiting_provider"`
	Error             int     `json:"error"`
	Complete          int     `json:"complete"`
	RefundTotal       int     `json:"refund_total"`
	RefundFee         int     `json:"refund_fee"`
	Refunded          int     `json:"refunded"`
	RefundedPartially int     `json:"refunded_partially"`
	Total             int     `json:"total"`
	Volume            float64 `json:"volume"`
	VolumeFee         float64 `json:"volume_fee"`
}
