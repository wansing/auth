package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/wansing/auth/server"
)

var db *server.Database

func insert() {

	var username, password, repeatPassword  string

	fmt.Print("Username: ")
	if _, err := fmt.Scanln(&username); err != nil {
		log.Println(err)
		return
	}

	fmt.Print("Password: ")
	if _, err := fmt.Scanln(&password); err != nil {
		log.Println(err)
		return
	}

	fmt.Print("Repeat password: ")
	if _, err := fmt.Scanln(&repeatPassword); err != nil {
		log.Println(err)
		return
	}

	if password == "" {
		log.Println("Password is empty")
		return
	}

	if password != repeatPassword {
		log.Println("Repetition doesn't match")
		return
	}

	if err := db.Insert(username, password); err != nil {
		log.Printf("Error inserting: %v", err)
		return
	}
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: authcli [flags] command\n  command\n        insert|TODO\n")
		flag.PrintDefaults()
	}

	dbDriver, dbDSN := server.DBFlags()
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Please specify a command")
		return
	}
	action := strings.ToLower(flag.Arg(0))

	// open database

	var err error
	db, err = server.OpenDatabase(*dbDriver, *dbDSN)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	log.Printf("Opened %s database %s", *dbDriver, *dbDSN)

	switch action {
	case "insert":
		insert()
	default:
		fmt.Printf("Unknown action: %s", action)
	}
}
