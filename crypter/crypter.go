package crypter

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func AcceptPassword() (string, error) {
	fmt.Print("Enter encryption password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func New(password string)
