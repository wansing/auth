package server

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"encoding/base64"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
	allStmt    *sql.Stmt
	authStmt   *sql.Stmt
	deleteStmt *sql.Stmt
	insertStmt *sql.Stmt
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
			name TEXT PRIMARY KEY, -- convention: lowercase
			salt TEXT NOT NULL,
			scheme TEXT NOT NULL,
			hash TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	db := &Database{DB: sqlDB}

	db.allStmt = db.mustPrepare("SELECT name FROM users")
	db.authStmt = db.mustPrepare("SELECT salt, scheme, hash FROM users WHERE name = ?")
	db.deleteStmt = db.mustPrepare("DELETE FROM users WHERE name = ?")
	db.insertStmt = db.mustPrepare("INSERT INTO users (name, salt, scheme, hash) VALUES (?, ?, ?, ?)")

	return db, nil
}

func SSHA512(password, salt string) string {
	sum := sha512.Sum512([]byte(password + salt)) // like dovecot's doveadm, we do just concatenate the plaintext password and the salt
	return base64.StdEncoding.EncodeToString(sum[:])                                 // convert array to slice
}

func (db *Database) Authenticate(username, password string) (success bool, err error) {

	rows, err := db.authStmt.Query(username)
	if err != nil {
		err = nil // user not found is not an error
		return
	}
	defer rows.Close()

	var salt, scheme string
	var dbHash string

	for rows.Next() {
		err = rows.Scan(&salt, &scheme, &dbHash)
		if err != nil {
			return
		}
		break
	}

	if rows.Next() {
		err = fmt.Errorf("Multiple results for user %s", username)
		return
	}

	if salt == "" || dbHash == "" {
		err = fmt.Errorf("Incomplete record for user %s", username)
		return
	}

	var passwordHash string
	switch strings.ToLower(scheme) {
	case "ssha512":
		passwordHash = SSHA512(password, salt)
	default:
		err = fmt.Errorf("Unknown hash scheme for user %s", username)
		return
	}

	success = (passwordHash == dbHash)
	return
}

func (db *Database) All() ([]string, error) {

	rows, err := db.allStmt.Query()
	if err != nil && err != sql.ErrNoRows {
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

	b := make([]byte, 24)

	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	salt := base64.StdEncoding.EncodeToString(b)

	_, err = db.insertStmt.Exec(username, salt, "ssha512", SSHA512(password, salt))
	return err
}
