package postgre

type Transaction struct {
	ID             string `json:"id"`
	CreditWalletID string `validate:"required" json:"creditWalletId"`
	DebitWalletID  string `validate:"required" json:"debitWalletId"`
	Amount         int    `validate:"required" json:"amount"`
	Type           int    `json:"type"`
	FeeAmount      int    `json:"feeAmount"`
	FeeWalletID    string `json:"feeWalletId"`
	CreditUserID   string `json:"creditUserId"`
	DebitUserID    string `json:"debitUserId"`
}
