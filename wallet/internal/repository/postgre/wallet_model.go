package postgre

type Wallet struct {
	ID      string `json:"id"`
	Balance int    `validate:"required" json:"balance"`
	UserID  string `validate:"required" json:"userId"`
}
