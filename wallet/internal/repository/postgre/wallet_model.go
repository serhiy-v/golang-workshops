package postgre

type Wallet struct {
	Id      string `json:"id"`
	Balance int    `validate:"required" json:"balance"`
	UserId  string `validate:"required" json:"userId"`
}
