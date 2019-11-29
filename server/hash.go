package server

import (
	"strings"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Scheme interface {
	Compare(hashedPassword []byte, password string) error
	Generate(password string) ([]byte, error)
	Name() string
}

type Bcrypt struct{}

func (_ Bcrypt) Compare(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func (_ Bcrypt) Generate(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (_ Bcrypt) Name() string {
	return "bcrypt"
}

func GetScheme(scheme string) (Scheme, error) {
	switch strings.ToLower(scheme) {
	case "bcrypt":
		return Bcrypt{}, nil
	default:
		return nil, fmt.Errorf("Unknown scheme: %s", scheme)
	}
}
