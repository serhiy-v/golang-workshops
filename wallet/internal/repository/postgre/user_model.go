package postgre

import "database/sql"

type User struct {
	Id    string         `json:"id"`
	Name  string         `json:"name"`
	Token sql.NullString `json:"token"`
}
