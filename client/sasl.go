package client

import (
	"io/ioutil"
	"net"
	"strings"

	"github.com/emersion/go-sasl"
)

// Authenticates against an unix socket which receivers a SASL PLAIN request and returns "authenticated" if authentication was successful.
type SASLPlain struct {
	Socket string
}

func (s SASLPlain) Authenticate(email, password string) (bool, error) {

	conn, err := net.Dial("unix", s.Socket)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	_, initialResponse, err := sasl.NewPlainClient("", email, password).Start()
	if err != nil {
		return false, err
	}

	_, err = conn.Write(initialResponse)
	if err != nil {
		return false, err
	}

	err = conn.(*net.UnixConn).CloseWrite()
	if err != nil {
		return false, err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return false, err
	}

	if strings.ToLower(string(result)) == "authenticated" {
		return true, nil
	} else {
		return false, nil
	}
}

func (s SASLPlain) Available() bool {
	return s.Socket != ""
}

func (s SASLPlain) Name() string {
	return "SASL Plain via Unix Socket (non-standardized)"
}
