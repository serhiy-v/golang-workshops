package postgre

type Transaction struct {
	Id             string `json:"id"`
	CreditWalletId string `validate:"required" json:"creditWalletId"`
	DebitWalletId  string `validate:"required" json:"debitWalletId"`
	Amount         int    `validate:"required" json:"amount"`
	Type           int    `json:"type"`
	FeeAmount      int    `json:"feeAmount"`
	FeeWalletId    string `json:"feeWalletId"`
	CreditUserId   string `json:"creditUserId"`
	DebitUserId    string `json:"debitUserId"`
}
