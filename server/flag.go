package server

import "flag"

func DBFlags() (db, dsn *string) {
	db = flag.String("db", "sqlite3", `database driver, can be "mysql" (untested), "postgres" (untested) or "sqlite3"`)
	dsn = flag.String("dsn", "auth.sqlite3", "database data source name")
	return
}
