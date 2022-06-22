package models

type User struct {
	ID    string  `json:"id" bson:"_id"`
	Name  string  `validate:"required" json:"name"`
	Token *string `json:"token" bson:"token"`
}
