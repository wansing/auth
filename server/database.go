package server

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"encoding/base64"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
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
			username TEXT PRIMARY KEY, -- convention: lowercase
			salt TEXT NOT NULL,
			scheme TEXT NOT NULL,
			hash TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	db := &Database{DB: sqlDB}

	db.authStmt = db.mustPrepare("SELECT salt, scheme, hash FROM users WHERE username = ?")
	db.deleteStmt = db.mustPrepare("DELETE FROM users WHERE username = ?")
	db.insertStmt = db.mustPrepare("INSERT INTO users (username, salt, scheme, hash) VALUES (?, ?, ?, ?)")

	return db, nil
}

func SSHA512(password, salt string) string {
	sum := sha512.Sum512([]byte(password + salt)) // like dovecot's doveadm, we do just concatenate the plaintext password and the salt
	return base64.StdEncoding.EncodeToString(sum[:])                                 // convert array to slice
}

func (s *Database) Authenticate(username, password string) (success bool, err error) {

	rows, err := s.authStmt.Query(username)
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

func (s *Database) Delete(username string) error {
	_, err := s.deleteStmt.Exec(username)
	return err
}

func (s *Database) Insert(username, password string) error {

	b := make([]byte, 24)

	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	salt := base64.StdEncoding.EncodeToString(b)

	_, err = s.insertStmt.Exec(username, salt, "ssha512", SSHA512(password, salt))
	return err
}
