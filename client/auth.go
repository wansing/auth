package client

import "log"

type Authenticators []Authenticator

type Authenticator interface {
	Authenticate(email, password string) (success bool, err error)
	Available() bool
	Name() string
}

func (as Authenticators) Authenticate(email, password string) (success bool, err error) {

	for _, a := range as {

		if !a.Available () {
			continue
		}

		success, err = a.Authenticate(email, password)
		if success {
			break
		}
		if err != nil {
			log.Printf("Error authenticating with %s authenticator: %v", a.Name(), err)
		}
	}
	return
}

func (as Authenticators) Available() bool {
	for _, a := range as {
		if a.Available() {
			return true
		}
	}
	return false
}