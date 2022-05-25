package postgre

import (
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
