package helpers

import (
	"log"
	"net/http"
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

func GenerateJwt(name string, email string) (string, string, error) {
	claims := &Info{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(5)).Unix(), //Create the access token (expires in 5 minutes)
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
		return "There is error in signing the token", "", err
	}
	log.Println("Signed token accesstoken is :", SignedAccessToken)
	//we created the refresh token here ... expire time is 24 hours now..
	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":   name,
		"email":  email,
		"expire": jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()},
	})
	SignedRefreshToken, err := RefreshToken.SignedString([]byte(Secret_key))
	if err != nil {
		return "There is error in signing the token", "", err
	}
	log.Println("Signed token refreshtoken is :", SignedRefreshToken)

	return SignedAccessToken, SignedRefreshToken, nil
}

//If we need implement the automatic generation of Accesstoken use GOROUTINE
// func init() {
// 	HandleRefreshToken(SignedRefreshToken)
// }

// Need to handle the refresh token and access token about When should the refresh token should create new access token and all
func HandleRefreshToken(refreshTokenString string) (*TokenString, error) {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Fetch the secret key from .env file
	SecretKey := os.Getenv("SECRET_KEY")

	claims := &Info{}
	log.Println(claims)
	//RefreshToken := jwt.Parse(helpers.SignedAccessToken, func(t *jwt.Token) (interface{}, error))
	RefreshToken, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(SecretKey), nil
	})
	if err != nil || !RefreshToken.Valid {
		return &TokenString{AccessToken: "Invalid refresh token and refresh also expired "}, http.ErrNoCookie
	}
	expireTime := time.Unix(claims.ExpiresAt, 0)
	log.Println(expireTime)
	remaingTime := time.Until(expireTime)
	if remaingTime > 30*time.Second {
		return &TokenString{AccessToken: "Token is not yet Expired..."}, nil
	}
	username := claims.Name
	email := claims.Email
	NewClaims := &Info{
		Name:  username,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(5)).Unix(), // 5 mins
		},
	}
	NewAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims)
	NewAccessTokenString, err := NewAccessToken.SignedString([]byte(SecretKey))
	if err != nil {
		log.Print("Failed to sign new access token:", err)
		return nil, err
	}
	return &TokenString{AccessToken: NewAccessTokenString}, nil
}

type TokenString struct {
	AccessToken string `json:"accesstoken"`
}
