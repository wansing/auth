package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/wansing/auth/server"
	"golang.org/x/crypto/ssh/terminal"
)

var db *server.Database

func all() {
	if all, err := db.All(); err == nil {
		for _, username := range all {
			fmt.Println(username)
		}
	} else {
		fmt.Println(err)
	}
}

func delete(username string) {
	if err := db.Delete(username); err != nil {
		fmt.Println(err)
		return
	}
}

func insert(username string) {

	var password, repeatPassword []byte
	var err error

	fmt.Print("Password: ")
	if password, err = terminal.ReadPassword(0); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()

	if len(password) == 0 {
		fmt.Println("Password is empty")
		return
	}

	fmt.Print("Repeat password: ")
	if repeatPassword, err = terminal.ReadPassword(0); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()

	if !bytes.Equal(password, repeatPassword) {
		fmt.Println("Repetition doesn't match")
		return
	}

	if err := db.Insert(username, string(password)); err != nil {
		fmt.Printf("Error inserting: %v\n", err)
		return
	}
}

func verify(username string) {

	var password []byte
	var err error

	fmt.Print("Password: ")
	if password, err = terminal.ReadPassword(0); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()

	if success, err := db.Authenticate(username, string(password)); err == nil {
		if success {
			fmt.Println("Verification ok")
		} else {
			fmt.Println("Verification failed")
		}
	} else {
		fmt.Printf("Error verifying: %v", err)
		return
	}
}

func main() {

	argAll := flag.Bool("all", false, "list all usernames")
	argDelete := flag.String("delete", "", "delete a user")
	argInsert := flag.String("insert", "", "insert a user")
	argVerify := flag.String("verify", "", "verify the password of a user")
	dbDriver, dbDSN := server.DBFlags()
	flag.Parse()

	// open database

	var err error
	db, err = server.OpenDatabase(*dbDriver, *dbDSN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	switch {
	case *argInsert != "":
		insert(*argInsert)
	case *argAll:
		all()
	case *argDelete != "":
		delete(*argDelete)
	case *argVerify != "":
		verify(*argVerify)
	default:
		fmt.Println("no command specified")
	}
}
