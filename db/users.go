package db

import (
	"database/sql"
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
