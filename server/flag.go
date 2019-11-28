package server

import "flag"

func DBFlags() (dbDriver, dbDSN *string) {
	dbDriver = flag.String("db", "sqlite3", `database driver, can be "mysql" (untested), "postgres" (untested) or "sqlite3"`)
	dbDSN = flag.String("dsn", "auth.sqlite3", "database data source name")
	return
}
