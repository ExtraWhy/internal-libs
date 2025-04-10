package db

import (
	"database/sql"
	"fmt"

	"github.com/ExtraWhy/internal-libs/models/user"
	"github.com/google/uuid"
)

func CreateUsersTable(db *sql.DB) error {
	const tableSQL = `CREATE TABLE users(
  id TEXT PRIMARY KEY,
  username TEXT,
  email TEXT,
  token TEXT,
  picture TEXT
);`

	statement, err := db.Prepare(tableSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec()
	return err
}

func (dbc *DBSqlConnection) InsertUser(u user.User) error {
	// Insert only if new user
	var count int
	err := dbc.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", u.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user with email %s already exists", u.Email)
	}

	u.Id = uuid.New().String()

	//If it's a new user proceed to create it
	query := `INSERT INTO users(id, username, email, token, picture) VALUES (?, ?, ?, ?, ?)`
	stmt, err := dbc.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Id, u.Username, u.Email, u.Token, u.Picture)

	return err
}

func (dbc *DBSqlConnection) GetUsers() ([]user.User, error) {
	rows, err := dbc.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Token, &u.Picture); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil

}
