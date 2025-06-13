package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPsw(Password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 13)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func VerifyPasw(Password string, HashedPassword string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(Password)); err != nil {
		return "Password is incorrect ", err
	}

	return "Password Matched ", nil
}
