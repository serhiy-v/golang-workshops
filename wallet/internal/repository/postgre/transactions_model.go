package postgre

type Transaction struct {
	Id             string `json:"id"`
	CreditWalletId string `json:"creditWalletId"`
	DebitWalletId  string `json:"debitWalletId"`
	Amount         int    `json:"amount"`
	FeeAmount      int    `json:"feeAmount"`
	FeeWalletId    string `json:"feeWalletId"`
	CreditUserId   string `json:"creditUserId"`
	DebitUserId    string `json:"debitUserId"`
}
