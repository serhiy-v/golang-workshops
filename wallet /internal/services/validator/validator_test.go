package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/workshops/wallet/internal/repository/models"
)

//nolint
func TestValidateUser(t *testing.T) {
	validate := NewValidator()

	user := &models.User{
		Name: "test",
	}

	err := validate.Validate(user)
	assert.NoError(t, err)
}

//nolint
func TestValidateUserError(t *testing.T) {
	validate := NewValidator()

	user := &models.User{
		ID: "test",
	}

	err := validate.Validate(user)
	assert.Error(t, err)
}

//nolint
func TestValidateWallet(t *testing.T) {
	validate := NewValidator()

	wallet := &models.Wallet{
		Balance: 100,
		UserID:  "92f0d2ea-f6ac-4b20-bb20-01062b29eb9a",
	}

	err := validate.Validate(wallet)
	assert.NoError(t, err)
}

//nolint
func TestValidateWalletError(t *testing.T) {
	validate := NewValidator()

	wallet := &models.Wallet{
		Balance: 100,
	}

	err := validate.Validate(wallet)
	assert.Error(t, err)
}

//nolint
func TestValidateTransaction(t *testing.T) {
	validate := NewValidator()

	transaction := &models.Transaction{
		CreditWalletID: "ce71eb21-1312-4e29-89df-039cae56007a",
		DebitWalletID:  "096a20c7-0b2a-475a-b175-229196f23cde",
		Amount:         20,
	}

	err := validate.Validate(transaction)
	assert.NoError(t, err)
}

//nolint
func TestValidateTransactionError(t *testing.T) {
	validate := NewValidator()

	transaction := &models.Transaction{
		CreditWalletID: "ce71eb21-1312-4e29-89df-039cae56007a",
		Amount:         20,
	}

	err := validate.Validate(transaction)
	assert.Error(t, err)
}
