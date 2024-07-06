package repository // доступ к базе данных и операции с ней

import (
    "github.com/jmoiron/sqlx"
	"fmt"
	"orderTracking"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user orderTracking.User) (int, error){
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err!= nil{
		return 0, err
	}
	return id, nil
}