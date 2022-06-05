package wallet

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/workshops/wallet/internal/repository/postgre"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	expectedUser := []*postgre.User{
		&postgre.User{
			Id:    "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
			Token: sql.NullString{"yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic2VyaGlpIiwiZXhwIjoxNjU3MTcxMjYxfQ.p9B8ZZFmYtF6euIdDQJA9NbeCJaGCUXHxMh8wR0VyWw", true},
		},
		&postgre.User{
			Id:    "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
			Token: sql.NullString{"yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic2VyaGlpIiwiZXhwIjoxNjU3MTcxMjYxfQ.p9B8ZZFmYtF6euIdDQJA9NbeCJaGCUXHxMh8wR0VyWw", true},
		},
	}

	q := "SELECT * FROM users"

	mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnRows(mock.NewRows([]string{"id", "token"}).AddRow("928eeecf-05ad-4e6f-ab7f-5477225b4c52", sql.NullString{"yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic2VyaGlpIiwiZXhwIjoxNjU3MTcxMjYxfQ.p9B8ZZFmYtF6euIdDQJA9NbeCJaGCUXHxMh8wR0VyWw", true}).AddRow("928eeecf-05ad-4e6f-ab7f-5477225b4c52", sql.NullString{"yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic2VyaGlpIiwiZXhwIjoxNjU3MTcxMjYxfQ.p9B8ZZFmYtF6euIdDQJA9NbeCJaGCUXHxMh8wR0VyWw", true}))

	user, err := srvc.GetUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

}

func TestGetUsersError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	mockErr := errors.New("Error getting users")

	q := "SELECT * FROM users"

	mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnError(mockErr)

	_, err = srvc.GetUsers()

	assert.Error(t, err)
	assert.Equal(t, err, mockErr)
}

func TestCreateUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	token := "yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic2VyaGlpIiwiZXhwIjoxNjU3MTcxMjYxfQ.p9B8ZZFmYtF6euIdDQJA9NbeCJaGCUXHxMh8wR0VyWw"

	q := "INSERT INTO users (token) VALUES ($1)"

	mock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(token).WillReturnResult(sqlmock.NewResult(1, 1))

	err = srvc.CreateUser(token)

	assert.NoError(t, err)
}

func TestCreateUsersError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	mockErr := errors.New("Unable to create users")

	token := "yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic2VyaGlpIiwiZXhwIjoxNjU3MTcxMjYxfQ.p9B8ZZFmYtF6euIdDQJA9NbeCJaGCUXHxMh8wR0VyWw"

	q := "INSERT INTO users (token) VALUES ($1)"

	mock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(token).WillReturnError(mockErr)

	err = srvc.CreateUser(token)

	assert.Error(t, err)
	assert.Equal(t, err, mockErr)
}

func TestCreateWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "INSERT INTO wallets (balance, user_id) VALUES ($1,$2)"

	wallet := &postgre.Wallet{
		Balance: 100,
		UserId:  "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
	}

	mock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(wallet.Balance, wallet.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

	err = srvc.CreateWallet(wallet)

	assert.NoError(t, err)
}

func TestCreateWalletError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "INSERT INTO wallets (balance, user_id) VALUES ($1,$2)"

	wallet := &postgre.Wallet{
		Balance: 100,
		UserId:  "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
	}

	mockErr := errors.New("Unable to create wallet")

	mock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(wallet.Balance, wallet.UserId).WillReturnError(mockErr)

	err = srvc.CreateWallet(wallet)

	assert.Error(t, err)
	assert.Equal(t, err, mockErr)
}

func TestGetWalletById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "SELECT id,balance,user_id FROM wallets WHERE id=$1"

	expectedWallet := &postgre.Wallet{
		Id:      "096a20c7-0b2a-475a-b175-229196f23cde",
		Balance: 100,
		UserId:  "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
	}

	id := "096a20c7-0b2a-475a-b175-229196f23cde"

	mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnRows(mock.NewRows([]string{"id", "balance", "userId"}).AddRow(expectedWallet.Id, expectedWallet.Balance, expectedWallet.UserId))

	wallet, err := srvc.GetWalletByID(id)

	assert.NoError(t, err)
	assert.Equal(t, wallet, expectedWallet)
}

func TestGetWalletByIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "SELECT id,balance,user_id FROM wallets WHERE id=$1"

	mockErr := errors.New("Unable to get wallet by id")

	id := "096a20c7-0b2a-475a-b175-229196f23cde"

	mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnError(mockErr)

	_, err = srvc.GetWalletByID(id)

	assert.Error(t, err)
	assert.Equal(t, err, mockErr)
}

func TestGetWalletTransactionsById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "SELECT * FROM transactions WHERE credit_wallet_id=$1 or debit_wallet_id=$1"

	id := "ce71eb21-1312-4e29-89df-039cae56007a"

	expectedTransaction := []*postgre.Transaction{
		&postgre.Transaction{
			Id:             "a15abc6c-63c5-46a4-bf0c-f355a23edc2e",
			CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
			DebitWalletId:  "096a20c7-0b2a-475a-b175-229196f23cde",
			Amount:         20,
			Type:           1,
			FeeAmount:      3,
			FeeWalletId:    "85aa7525-4fdb-4436-a600-66ffc55e0f65",
			CreditUserId:   "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
			DebitUserId:    "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
		},
		&postgre.Transaction{
			Id:             "a15abc6c-63c5-46a4-bf0c-f355a23edc2e",
			CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
			DebitWalletId:  "096a20c7-0b2a-475a-b175-229196f23cde",
			Amount:         20,
			Type:           1,
			FeeAmount:      3,
			FeeWalletId:    "85aa7525-4fdb-4436-a600-66ffc55e0f65",
			CreditUserId:   "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
			DebitUserId:    "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
		},
	}

	mock.ExpectQuery(regexp.QuoteMeta(q)).WithArgs(id).WillReturnRows(mock.NewRows([]string{"id", "creditWalletId", "debitWalletId", "amount", "type", "feeAmount", "feeWalletId", "creditUserId", "debitUserId"}).AddRow("a15abc6c-63c5-46a4-bf0c-f355a23edc2e", "ce71eb21-1312-4e29-89df-039cae56007a", "096a20c7-0b2a-475a-b175-229196f23cde", 20, 1, 3, "85aa7525-4fdb-4436-a600-66ffc55e0f65", "928eeecf-05ad-4e6f-ab7f-5477225b4c52", "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a").AddRow("a15abc6c-63c5-46a4-bf0c-f355a23edc2e", "ce71eb21-1312-4e29-89df-039cae56007a", "096a20c7-0b2a-475a-b175-229196f23cde", 20, 1, 3, "85aa7525-4fdb-4436-a600-66ffc55e0f65", "928eeecf-05ad-4e6f-ab7f-5477225b4c52", "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a"))

	transaction, err := srvc.GetWalletTransactionsById(id)

	assert.NoError(t, err)
	assert.Equal(t, expectedTransaction, transaction)
}

func TestGetWalletTransactionsByIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "SELECT * FROM transactions WHERE credit_wallet_id=$1 or debit_wallet_id=$1"

	id := "ce71eb21-1312-4e29-89df-039cae56007a"

	mockErr := errors.New("Unable to get transaction")

	mock.ExpectQuery(regexp.QuoteMeta(q)).WithArgs(id).WillReturnError(mockErr)

	_, err = srvc.GetWalletTransactionsById(id)

	assert.Error(t, err)
	assert.Equal(t, err, mockErr)
}

func TestGetransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "SELECT * FROM transactions"

	expectedTransaction := []*postgre.Transaction{
		&postgre.Transaction{
			Id:             "a15abc6c-63c5-46a4-bf0c-f355a23edc2e",
			CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
			DebitWalletId:  "096a20c7-0b2a-475a-b175-229196f23cde",
			Amount:         20,
			Type:           1,
			FeeAmount:      3,
			FeeWalletId:    "85aa7525-4fdb-4436-a600-66ffc55e0f65",
			CreditUserId:   "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
			DebitUserId:    "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
		},
		&postgre.Transaction{
			Id:             "a15abc6c-63c5-46a4-bf0c-f355a23edc2e",
			CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
			DebitWalletId:  "096a20c7-0b2a-475a-b175-229196f23cde",
			Amount:         20,
			Type:           1,
			FeeAmount:      3,
			FeeWalletId:    "85aa7525-4fdb-4436-a600-66ffc55e0f65",
			CreditUserId:   "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
			DebitUserId:    "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
		},
	}

	mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnRows(mock.NewRows([]string{"id", "creditWalletId", "debitWalletId", "amount", "type", "feeAmount", "feeWalletId", "creditUserId", "debitUserId"}).AddRow("a15abc6c-63c5-46a4-bf0c-f355a23edc2e", "ce71eb21-1312-4e29-89df-039cae56007a", "096a20c7-0b2a-475a-b175-229196f23cde", 20, 1, 3, "85aa7525-4fdb-4436-a600-66ffc55e0f65", "928eeecf-05ad-4e6f-ab7f-5477225b4c52", "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a").AddRow("a15abc6c-63c5-46a4-bf0c-f355a23edc2e", "ce71eb21-1312-4e29-89df-039cae56007a", "096a20c7-0b2a-475a-b175-229196f23cde", 20, 1, 3, "85aa7525-4fdb-4436-a600-66ffc55e0f65", "928eeecf-05ad-4e6f-ab7f-5477225b4c52", "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a"))

	transaction, err := srvc.GetTransactions()

	assert.NoError(t, err)
	assert.Equal(t, expectedTransaction, transaction)
}

func TestGetTransactionsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)

	srvc := NewService(repo)

	q := "SELECT * FROM transactions"

	mockErr := errors.New("Unable to get transaction")

	mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnError(mockErr)

	_, err = srvc.GetTransactions()

	assert.Error(t, err)
	assert.Equal(t, err, mockErr)
}

//func TestCreateTransactions(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatal("Unable to connect")
//	}
//	defer db.Close()
//	repo := postgre.NewRepository(db)
//
//	srvc := NewService(repo)
//
//	q1 := "UPDATE wallets SET balance=balance-$1 WHERE id=$2"
//
//	q2 := "UPDATE wallets SET balance=balance+$1 WHERE id=$2"
//
//	q3 := "UPDATE wallets SET balance=balance+$1 WHERE id=$2"
//
//	q4 := "INSERT INTO transactions (credit_wallet_id,debit_wallet_id,amount,type,fee_amount,fee_wallet_id,credit_user_id, debit_user_id) VALUES ($1,$2,$3,$4,$5,$6,(SELECT user_id FROM wallets WHERE id=$7),(SELECT user_id FROM wallets WHERE id=$8))"
//
//	expectedTransaction := &postgre.Transaction{
//		Id:             "a15abc6c-63c5-46a4-bf0c-f355a23edc2e",
//		CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
//		DebitWalletId:  "096a20c7-0b2a-475a-b175-229196f23cde",
//		Amount:         20,
//		Type:           1,
//		FeeAmount:      3,
//		FeeWalletId:    "85aa7525-4fdb-4436-a600-66ffc55e0f65",
//		CreditUserId:   "928eeecf-05ad-4e6f-ab7f-5477225b4c52",
//		DebitUserId:    "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
//	}
//
//	mock.ExpectBegin()
//	mock.ExpectExec(regexp.QuoteMeta(q1)).WithArgs(expectedTransaction.Amount+2, expectedTransaction.CreditWalletId)
//	mock.ExpectExec(regexp.QuoteMeta(q2)).WithArgs(expectedTransaction.Amount, expectedTransaction.DebitWalletId)
//	mock.ExpectExec(regexp.QuoteMeta(q3)).WithArgs(expectedTransaction.FeeAmount, expectedTransaction.FeeWalletId)
//	mock.ExpectExec(regexp.QuoteMeta(q4)).WithArgs(expectedTransaction.CreditWalletId, expectedTransaction.DebitWalletId, expectedTransaction.Amount, expectedTransaction.Type, expectedTransaction.FeeAmount, expectedTransaction.FeeWalletId, expectedTransaction.CreditUserId, expectedTransaction.DebitUserId)
//	mock.ExpectClose()
//
//	err = srvc.CreateTransaction(expectedTransaction)
//
//	assert.NoError(t, err)
//}
