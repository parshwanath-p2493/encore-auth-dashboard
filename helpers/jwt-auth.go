package helpers

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Info struct {
	Name  string
	Email string
	jwt.StandardClaims
}

func GenerateJwt(name string, email string) (string, error) {
	claims := &Info{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(24)).Unix(), //Create the access token (expires in 15 minutes)
		},
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Secret_key := os.Getenv("SECRET_KEY")
	log.Println("Key is :" + Secret_key)
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println(AccessToken)
	SignedAccessToken, err := AccessToken.SignedString([]byte(Secret_key))
	if err != nil {
		return "There is error in signing the token", err
	}
	log.Println("Signed token accesstoken is :", SignedAccessToken)
	return SignedAccessToken, nil
}
