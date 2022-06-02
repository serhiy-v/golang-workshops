package postgre

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	Conn *sql.DB
}

func NewRepository(dsn string) (*Repository, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{Conn: conn}, nil
}

func (r *Repository) CreateUser(token string) error {
	q := "INSERT INTO users (token) VALUES ($1)"
	_, err := r.Conn.Exec(q, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUsers() ([]*User, error) {
	rows, err := r.Conn.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)

	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Id, &user.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateWallet(wallet *Wallet) error {
	q := "INSERT INTO wallets (balance, user_id) VALUES ($1,$2)"
	_, err := r.Conn.Exec(q, wallet.Balance, wallet.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetWalletByID(id string) (*Wallet, error) {
	q := "SELECT id,balance,user_id FROM wallets WHERE id=$1"
	wallet := new(Wallet)
	err := r.Conn.QueryRow(q, id).Scan(&wallet.Id, &wallet.Balance, &wallet.UserId)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *Repository) GetWalletTransactionsById(id string) ([]*Transaction, error) {
	rows, err := r.Conn.Query("SELECT * FROM transactions WHERE credit_wallet_id=$1 or debit_wallet_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]*Transaction, 0)

	for rows.Next() {
		transaction := new(Transaction)
		err := rows.Scan(&transaction.Id, &transaction.CreditWalletId, &transaction.DebitWalletId, &transaction.Amount, &transaction.Type, &transaction.FeeAmount, &transaction.FeeWalletId, &transaction.CreditUserId, &transaction.DebitUserId)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *Repository) GetTransactions() ([]*Transaction, error) {
	rows, err := r.Conn.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]*Transaction, 0)

	for rows.Next() {
		transaction := new(Transaction)
		err := rows.Scan(&transaction.Id, &transaction.CreditWalletId, &transaction.DebitWalletId, &transaction.Amount, &transaction.Type, &transaction.FeeAmount, &transaction.FeeWalletId, &transaction.CreditUserId, &transaction.DebitUserId)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *Repository) CreateTransaction(transaction *Transaction) error {
	ctx := context.Background()

	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=balance-$1 WHERE id=$2", transaction.Amount+2, transaction.CreditWalletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=balance+$1 WHERE id=$2", transaction.Amount, transaction.DebitWalletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=balance+$1 WHERE id=$2", 2, "85aa7525-4fdb-4436-a600-66ffc55e0f65")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (credit_wallet_id,debit_wallet_id,amount,type,fee_amount,fee_wallet_id,credit_user_id, debit_user_id) VALUES ($1,$2,$3,$4,$5,$6,(SELECT user_id FROM wallets WHERE id=$7),(SELECT user_id FROM wallets WHERE id=$8))", transaction.CreditWalletId, transaction.DebitWalletId, transaction.Amount, 1, 2, "85aa7525-4fdb-4436-a600-66ffc55e0f65", transaction.CreditWalletId, transaction.DebitWalletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
