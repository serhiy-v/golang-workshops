package postgre

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/workshops/wallet/internal/repository/models"
)

type Repository struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{Conn: conn}
}

func NewPostgresDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return conn, nil
}

func (r *Repository) CreateUser(token string) error {
	q := "INSERT INTO users (token) VALUES ($1)"
	_, err := r.Conn.Exec(q, token)

	if err != nil {
		return errors.Wrap(err, "Error from db")
	}

	return nil
}

func (r *Repository) GetUsers() ([]*models.User, error) {
	rows, err := r.Conn.Query("SELECT * FROM users")
	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	defer rows.Close()

	users := make([]*models.User, 0)

	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.ID, &user.Token)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return users, nil
}

func (r *Repository) CreateWallet(wallet *models.Wallet) error {
	q := "INSERT INTO wallets (balance, user_id) VALUES ($1,$2)"
	_, err := r.Conn.Exec(q, wallet.Balance, wallet.UserID)

	if err != nil {
		return errors.Wrap(err, "Error from db")
	}

	return nil
}

func (r *Repository) GetWalletByID(id string) (*models.Wallet, error) {
	q := "SELECT id,balance,user_id FROM wallets WHERE id=$1"
	wallet := new(models.Wallet)
	err := r.Conn.QueryRow(q, id).Scan(&wallet.ID, &wallet.Balance, &wallet.UserID)

	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return wallet, nil
}

func (r *Repository) GetWalletTransactionsByID(id string) ([]*models.Transaction, error) {
	rows, err := r.Conn.Query("SELECT * FROM transactions WHERE credit_wallet_id=$1 or debit_wallet_id=$1", id)
	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}
	defer rows.Close()

	transactions := make([]*models.Transaction, 0)

	for rows.Next() {
		transaction := new(models.Transaction)
		err := rows.Scan(&transaction.ID, &transaction.CreditWalletID, &transaction.DebitWalletID, &transaction.Amount,
			&transaction.Type, &transaction.FeeAmount, &transaction.FeeWalletID,
			&transaction.CreditUserID, &transaction.DebitUserID)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return transactions, nil
}

func (r *Repository) GetTransactions() ([]*models.Transaction, error) {
	rows, err := r.Conn.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	defer rows.Close()

	transactions := make([]*models.Transaction, 0)

	for rows.Next() {
		transaction := new(models.Transaction)
		err := rows.Scan(&transaction.ID, &transaction.CreditWalletID, &transaction.DebitWalletID, &transaction.Amount,
			&transaction.Type, &transaction.FeeAmount, &transaction.FeeWalletID,
			&transaction.CreditUserID, &transaction.DebitUserID)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return transactions, nil
}

func (r *Repository) CreateTransaction(transaction *models.Transaction) error {
	ctx := context.Background()

	tx, err := r.Conn.BeginTx(ctx, nil)

	if err != nil {
		return errors.Wrap(err, "Error from db")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=balance-$1 WHERE id=$2",
		transaction.Amount+2, transaction.CreditWalletID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			log.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}

		return errors.Wrap(err, "Error from db")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=balance+$1 WHERE id=$2",
		transaction.Amount, transaction.DebitWalletID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			log.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}

		return errors.Wrap(err, "Error from db")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=balance+$1 WHERE id=$2",
		2, "85aa7525-4fdb-4436-a600-66ffc55e0f65")
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			log.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}

		return errors.Wrap(err, "Error from db")
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (credit_wallet_id,debit_wallet_id,amount,"+
		"type,fee_amount,fee_wallet_id,credit_user_id, debit_user_id,date) VALUES "+
		"($1,$2,$3,$4,$5,$6,(SELECT user_id FROM wallets WHERE id=$7),(SELECT user_id FROM wallets WHERE id=$8),"+
		"$9)",
		transaction.CreditWalletID, transaction.DebitWalletID, transaction.Amount, 1, 2,
		"85aa7525-4fdb-4436-a600-66ffc55e0f65", transaction.CreditWalletID, transaction.DebitWalletID, time.Now())
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			log.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}

		return errors.Wrap(err, "Error from db")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "Error from db")
	}

	return nil
}

func (r *Repository) GetWalletAmountDayByID(id string, week models.Week) ([]*models.Day, error) {
	q1 := "SELECT date_trunc('day',date) as day, SUM(amount) FROM transactions WHERE credit_wallet_id=$1 " +
		"AND date >= $2 AND date <= $3 GROUP BY date_trunc('day',date)"
	//badValue := 0
	days := make([]*models.Day, 0)
	rows, err := r.Conn.Query(q1, id, week.DateFrom, week.DateTo)

	for rows.Next() {
		day := new(models.Day)
		err := rows.Scan(&day.Date, &day.Outcome)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	q2 := "SELECT date_trunc('day',date), SUM(amount) FROM transactions WHERE debit_wallet_id=$1 " +
		"AND date >= $2 AND date <= $3 GROUP BY date_trunc('day',date)"

	rows, err = r.Conn.Query(q2, id, week.DateFrom, week.DateTo)
	for rows.Next() {
		incday := new(models.Day)
		err := rows.Scan(&incday.Date, &incday.Income)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		for _, day := range days {
			if day.Date == incday.Date {
				day.Income = incday.Income
			}
		}
	}

	return days, nil
}

func (r *Repository) GetWalletAmountWeekByID(id string, week models.Week) ([]*models.Day, error) {
	q1 := "SELECT date_trunc('week',date) as week,sum(amount) FROM transactions WHERE credit_wallet_id=$1 " +
		"AND date >= $2 AND date <= $3 GROUP BY date_trunc('week',date)"
	//badValue := 0
	days := make([]*models.Day, 0)
	rows, err := r.Conn.Query(q1, id, week.DateFrom, week.DateTo)

	for rows.Next() {
		day := new(models.Day)
		err := rows.Scan(&day.Date, &day.Outcome)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	q2 := "SELECT date_trunc('week',date) as week,sum(amount) FROM transactions WHERE debit_wallet_id=$1 " +
		"AND date >= $2 AND date <= $3 GROUP BY date_trunc('week',date)"

	rows, err = r.Conn.Query(q2, id, week.DateFrom, week.DateTo)
	for rows.Next() {
		incday := new(models.Day)
		err := rows.Scan(&incday.Date, &incday.Income)

		if err != nil {
			return nil, errors.Wrap(err, "Error from db")
		}

		for _, day := range days {
			if day.Date == incday.Date {
				day.Income = incday.Income
			}
		}
	}

	return days, nil
}
