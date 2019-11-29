package main

import (
	"bytes"
	"errors"
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

func promptPasswords() (string, error) {

	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Println()

	if len(password) == 0 {
		return "", errors.New("Password is empty")
	}

	fmt.Print("Repeat password: ")
	repeatPassword, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Println()

	if !bytes.Equal(password, repeatPassword) {
		return "", errors.New("Repetition doesn't match")
	}

	return string(password), nil
}

func insert(username string) {

	password, err := promptPasswords()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.Insert(username, password); err != nil {
		fmt.Printf("Error inserting: %v\n", err)
		return
	}
}

func reset(username string) {

	password, err := promptPasswords()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.Reset(username, password); err != nil {
		fmt.Printf("Error resetting: %v\n", err)
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

	if err := db.Authenticate(username, string(password)); err == nil {
		fmt.Println("Verification ok")
	} else {
		fmt.Printf("Verification error: %v\n", err)
		return
	}
}

func main() {

	argAll := flag.Bool("all", false, "list all usernames")
	argDelete := flag.String("delete", "", "delete a user")
	argInsert := flag.String("insert", "", "insert a user")
	argReset := flag.String("reset", "", "reset the password of a user")
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
	case *argReset != "":
		reset(*argReset)
	case *argVerify != "":
		verify(*argVerify)
	default:
		fmt.Println("no command specified")
	}
}
