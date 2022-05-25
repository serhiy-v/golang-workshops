package postgre

type Wallet struct {
	Id      string `json:"id"`
	Balance int    `json:"balance"`
	UserId  string `json:"userId"`
}
