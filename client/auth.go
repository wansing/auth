package client

import "fmt"

type Authenticators []Authenticator

type Authenticator interface {
	Authenticate(email, password string) (success bool, err error) // should not be called if Available() returns false
	Available() bool
	Name() string
}

// Tries the contained authenticators with the given credentials.
// If an authenticator returns an error, we continue with the next one. The last error is returned.
func (as Authenticators) Authenticate(email, password string) (bool, error) {

	var lastErr error

	for _, a := range as {

		if !a.Available () {
			continue
		}

		success, err := a.Authenticate(email, password)

		if err != nil {
			lastErr = fmt.Errorf("error querying authenticator %s: %v", a.Name(), err)
			success = false // security measure: in case of error, discard success
		}

		if success {
			return true, lastErr
		}
	}

	return false, lastErr
}

func (as Authenticators) Available() bool {
	for _, a := range as {
		if a.Available() {
			return true
		}
	}
	return false
}
