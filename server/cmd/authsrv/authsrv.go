package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/emersion/go-sasl"
	"github.com/wansing/auth/server"
)

var Db *server.Database
var wg sync.WaitGroup

func handle(conn net.Conn) {

	wg.Add(1)
	defer wg.Done()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Printf("error reading from connection: ", err)
		return
	}

	var authenticated = false

	saslServer := sasl.NewPlainServer(func(identity, username, password string) error {
		if err := Db.Authenticate(username, password); err == nil {
			authenticated = true
		}
		// we ignore other errors at the moment
		return nil // err would be returned by saslServer.Next()
	})

	_, done, _ := saslServer.Next(data) // err would come from our authenticator, which always returns nil

	if done && authenticated {
		_, err = conn.Write([]byte("authenticated"))
		if err != nil {
			log.Printf("error writing to connection: %v", err)
		}
	}

	_ = conn.Close()
}

func main() {

	log.SetFlags(0) // no log prefixes required, systemd-journald adds them

	dbDriver, dbDSN := server.DBFlags()
	socketPath := flag.String("socket", "./sasl.sock", "path of the SASL socket on which the application listens")
	flag.Parse()

	// open database

	var err error
	Db, err = server.OpenDatabase(*dbDriver, *dbDSN)
	if err != nil {
		log.Printf("error opening database: %v", err)
		return
	}
	defer Db.Close()

	log.Printf("opened %s database %s", *dbDriver, *dbDSN)

	// listen to socket

	fileinfo, err := os.Lstat(*socketPath)
	if err == nil {
		if fileinfo.Mode()&os.ModeSocket != 0 {
			_ = os.Remove(*socketPath)
		}
	}

	listener, err := net.Listen("unix", *socketPath)
	if err != nil {
		log.Printf("error creating socket: %v", err)
		return
	}
	defer listener.Close() // removes the socket file

	_ = os.Chmod(*socketPath, os.ModePerm) // chmod 777, so people can connect to the listener

	log.Printf("listening to %s", *socketPath)

	// accept connections

	var listenerClosed = false

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if listenerClosed {
					break
				} else {
					log.Printf("error accepting connection: %v", err)
					continue
				}
			}
			go handle(conn)
		}
	}()

	// wait for graceful shutdown

	sigintChannel := make(chan os.Signal, 1)
	signal.Notify(sigintChannel, os.Interrupt, syscall.SIGTERM) // SIGINT (Interrupt) or SIGTERM
	<-sigintChannel

	log.Println("received shutdown signal")
	listenerClosed = true
	listener.Close()
	wg.Wait()
}
