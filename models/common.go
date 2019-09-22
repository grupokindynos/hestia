package models

type BodyReq struct {
	Payload string `bson:"payload" json:"payload"`
}
