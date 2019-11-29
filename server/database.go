package server

import (
	"database/sql"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var preferredScheme = Bcrypt{}

type Database struct {
	*sql.DB
	allStmt    *sql.Stmt
	authStmt   *sql.Stmt
	deleteStmt *sql.Stmt
	insertStmt *sql.Stmt
	resetStmt  *sql.Stmt
}

func (db *Database) mustPrepare(query string) *sql.Stmt {
	if stmt, err := db.Prepare(query); err != nil {
		panic(err)
	} else {
		return stmt
	}
}

func OpenDatabase(backend, connStr string) (*Database, error) {

	sqlDB, err := sql.Open(backend, connStr)
	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			name TEXT PRIMARY KEY,
			scheme TEXT NOT NULL,
			hash BLOB NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	db := &Database{DB: sqlDB}

	db.allStmt = db.mustPrepare("SELECT name FROM users")
	db.authStmt = db.mustPrepare("SELECT scheme, hash FROM users WHERE name = ?")
	db.deleteStmt = db.mustPrepare("DELETE FROM users WHERE name = ?")
	db.insertStmt = db.mustPrepare("INSERT INTO users (name, scheme, hash) VALUES (?, ?, ?)")
	db.resetStmt = db.mustPrepare("UPDATE users SET scheme = ?, hash = ? WHERE name = ?")

	return db, nil
}

func (db *Database) Authenticate(username, password string) error {

	row := db.authStmt.QueryRow(username)

	var schemeStr string
	var hash []byte
	if err := row.Scan(&schemeStr, &hash); err != nil {
		return err
	}

	scheme, err := GetScheme(schemeStr)
	if err != nil {
		return err
	}

	return scheme.Compare(hash, password)
}

func (db *Database) All() ([]string, error) {

	rows, err := db.allStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	all := []string{}
	for rows.Next() {
		var username string
		rows.Scan(&username)
		all = append(all, username)
	}

	sort.Strings(all)

	return all, nil
}

func (db *Database) Delete(username string) error {
	_, err := db.deleteStmt.Exec(username)
	return err
}

func (db *Database) Insert(username, password string) error {

	hash, err := preferredScheme.Generate(password)
	if err != nil {
		return err
	}

	_, err = db.insertStmt.Exec(username, preferredScheme.Name(), hash)
	return err
}

func (db *Database) Reset(username, password string) error {

	hash, err := preferredScheme.Generate(password)
	if err != nil {
		return err
	}

	_, err = db.resetStmt.Exec(preferredScheme.Name(), hash, username)
	return err
}
