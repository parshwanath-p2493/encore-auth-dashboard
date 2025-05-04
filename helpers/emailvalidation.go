package helpers

import (
	"fmt"
	"net/mail"
	"strings"
)

func EmailValidation(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	Email, err := mail.ParseAddress(email)
	if err != nil {
		return "The email is not valid so give valid one", err
	}
	return fmt.Sprintf("The mail is unique and correct %v", Email), nil
}
