package models

type Payment struct {
	Address       string `bson:"address" json:"address"`
	Amount        string `bson:"amount" json:"amount"`
	Coin          string `bson:"coin" json:"coin"`
	Txid          string `bson:"txid" json:"txid"`
	Confirmations string `bson:"confirmations" json:"confirmations"`
}

type BodyReq struct {
	Payload string `bson:"payload" json:"payload"`
}
