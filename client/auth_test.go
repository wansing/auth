package client

import (
	"errors"
	"testing"
)

type Failer struct{}

func (Failer) Authenticate(email, password string) (success bool, err error) {
	return false, errors.New("error")
}

func (Failer) Name() string {
	return "Failer"
}

type AvailableFailer struct {
	Failer
}

func (AvailableFailer) Available() bool {
	return true
}

type UnavailableFailer struct {
	Failer
}

func (UnavailableFailer) Available() bool {
	return false
}

func TestAvailable(t *testing.T) {

	var as Authenticators
	as = append(as, AvailableFailer{})

	var got = as.Available()
	var want = true
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestUnavailable(t *testing.T) {

	var as Authenticators
	as = append(as, UnavailableFailer{})

	var got = as.Available()
	var want = false
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAuthenticate(t *testing.T) {

	var as Authenticators
	as = append(as, &SASLPlain{"/tmp/nonexisting-socket"}, SMTPS{1234}, STARTTLS{1234}, AvailableFailer{})

	gotSuccess, gotErr := as.Authenticate("some user", "some pass")
	wantSuccess := false
	wantErr := "error querying authenticator Failer: error"

	if gotSuccess != wantSuccess {
		t.Fatalf("got %v, want %v", gotSuccess, wantSuccess)
	}

	if gotErr.Error() != wantErr {
		t.Fatalf("got %v, want %v", gotErr, wantErr)
	}
}
