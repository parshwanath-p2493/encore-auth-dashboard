package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Info struct {
	Name  string
	Email string
}

func GenerateJwt(name string, email string) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")
	fmt.Println("Secret key", SECRET_KEY)
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //.SignedString([]byte(SECRET_KEY))
	fmt.Printf("\n SECRET_KEY: %s \n ", SECRET_KEY)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Print("Error signing the token", err)
		return " --->>", err
	}
	fmt.Printf("signedstring: %s", signedToken)
	return signedToken, nil
}
