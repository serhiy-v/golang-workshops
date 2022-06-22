package models

type Wallet struct {
	ID      string `json:"id" bson:"_id"`
	Balance int    `validate:"required" json:"balance" bson:"balance"`
	UserID  string `validate:"required" json:"userId" bson:"user_id"`
}
