package db

import (
	"database/sql"
	"fmt"

	"github.com/ExtraWhy/internal-libs/models/user"
)

func CreateUsersTable(db *sql.DB) error {
	const tableSQL = `CREATE TABLE users(
  id INTEGER PRIMARY KEY,
  username TEXT,
  email TEXT
);`

	statement, err := db.Prepare(tableSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec()
	return err
}

func (dbc *DBConnection) InsertUser(u user.User) error {
	// Insert only if new user
	var count int
	err := dbc.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", u.Id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user with id %d already exists", u.Id)
	}

	//If it's a new user proceed to create it
	query := `INSERT INTO users(id, username, email) VALUES (?, ?, ?)`
	stmt, err := dbc.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Id, u.Username, u.Email)
	return err
}

func (dbc *DBConnection) GetUsers() ([]user.User, error) {
	rows, err := dbc.db.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil

}
