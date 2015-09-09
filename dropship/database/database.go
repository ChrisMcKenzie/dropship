package database

import (
	"database/sql"

	_ "github.com/mxk/go-sqlite/sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "dropship.db")
	if err != nil {
		panic(err)
	}

	db.QueryRow(
		"create table if not exists tokens (repo STRING, token STRING)",
	)
}

func GetTokenFor(repo string) (res string, err error) {
	err = db.QueryRow(
		"SELECT token FROM tokens WHERE repo=?",
		repo,
	).Scan(&res)

	return
}

func StoreTokenFor(repo, token string) error {
	_, err := db.Exec(
		"INSERT INTO tokens (repo, token) VALUES (?, ?)",
		repo,
		token,
	)

	return err
}
