package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/workshops/wallet/internal/repository/postgre"
)

func TestValidateUser(t *testing.T) {
	validate := NewValidator()

	user := &postgre.User{
		Name: "test",
	}

	err := validate.Validate(user)
	assert.NoError(t, err)
}

func TestValidateUserError(t *testing.T) {
	validate := NewValidator()

	user := &postgre.User{
		Id: "test",
	}

	err := validate.Validate(user)
	assert.Error(t, err)
}

func TestValidateWallet(t *testing.T) {
	validate := NewValidator()

	wallet := &postgre.Wallet{
		Balance: 100,
		UserId:  "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
	}

	err := validate.Validate(wallet)
	assert.NoError(t, err)
}

func TestValidateWalletError(t *testing.T) {
	validate := NewValidator()

	wallet := &postgre.Wallet{
		Balance: 100,
	}

	err := validate.Validate(wallet)
	assert.Error(t, err)
}

func TestValidateTransaction(t *testing.T) {
	validate := NewValidator()

	transaction := &postgre.Transaction{
		CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
		DebitWalletId:  "096a20c7-0b2a-475a-b175-229196f23cde",
		Amount:         20,
	}

	err := validate.Validate(transaction)
	assert.NoError(t, err)
}

func TestValidateTransactionError(t *testing.T) {
	validate := NewValidator()

	transaction := &postgre.Transaction{
		CreditWalletId: "ce71eb21-1312-4e29-89df-039cae56007a",
		Amount:         20,
	}

	err := validate.Validate(transaction)
	assert.Error(t, err)
}
