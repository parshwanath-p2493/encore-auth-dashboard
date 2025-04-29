package helpers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"encore.app/database"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	GlobalSignedAccessToken  string
	GlobalSignedRefreshToken string
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
	//log.Println(AccessToken)
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
	GlobalSignedAccessToken = SignedAccessToken
	GlobalSignedRefreshToken = SignedRefreshToken

	return SignedAccessToken, SignedRefreshToken, nil
}

//If we need implement the automatic generation of Accesstoken use GOROUTINE
// func init() {
// 	HandleRefreshToken(SignedRefreshToken)
// }

// Need to handle the refresh token and access token about When should the refresh token should create new access token and all
func HandleRefreshToken() (*TokenString, error) {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("cursor is in HandleRefresh Token")
	//fmt.Println("cursor is in HandleRefresh Token")

	// Fetch the secret key from .env file
	SecretKey := os.Getenv("SECRET_KEY")
	AccessClaims := &Info{}
	AccessToken, err := jwt.ParseWithClaims(GlobalSignedAccessToken, AccessClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(SecretKey), nil
	})
	if err != nil || !AccessToken.Valid {
		log.Println("Invalid and access token  also expired LOgin again ")
	} else {

		expireTime := time.Unix(AccessClaims.ExpiresAt, 0)
		//log.Println("expire time of oldaccess token is", expireTime)
		remaingTime := time.Until(expireTime)
		if remaingTime > 30*time.Second {
			log.Println("Token is not yet Expired...")
			return &TokenString{AccessToken: "Token is not yet Expired..."}, nil
		}
		log.Println("Current time:", time.Now())
		log.Println("Token expiry time:", expireTime)
		log.Println("Remaining time:", remaingTime)
	}
	//Paese Refresh token for claims.....
	Refreshclaims := &Info{}
	log.Println(Refreshclaims)
	//RefreshToken := jwt.Parse(helpers.SignedAccessToken, func(t *jwt.Token) (interface{}, error))
	RefreshToken, err := jwt.ParseWithClaims(GlobalSignedRefreshToken, Refreshclaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(SecretKey), nil
	})
	if err != nil || !RefreshToken.Valid {
		return &TokenString{AccessToken: "Invalid refresh token and refresh also expired LOgin again "}, http.ErrNoCookie
	}
	log.Println("Cursorsor moved next step form access to refresh")
	log.Println(Refreshclaims.Name)
	log.Println(Refreshclaims.Email)

	username := Refreshclaims.Name
	email := Refreshclaims.Email
	NewClaims := &Info{
		Name:  username,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(5)).Unix(), // 5 mins
		},
	}
	NewAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims)
	NewAccessTokenString, err := NewAccessToken.SignedString([]byte(SecretKey))
	GlobalSignedAccessToken = NewAccessTokenString
	if err != nil {
		log.Print("Failed to sign new access token:", err)
		return nil, err
	}
	RefreshExpiretime := time.Unix(Refreshclaims.ExpiresAt, 0)
	log.Println("Refresh Token expiry time", RefreshExpiretime)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.OpenCollection("Users")
	filter := bson.M{"email": Refreshclaims.Email}
	update := bson.M{
		"$set": bson.M{
			"token":                   NewAccessTokenString,
			"expiretimetoken":         time.Now(),
			"expiretimerefresh_token": RefreshExpiretime,
			"refresh_token":           GlobalSignedRefreshToken,
			"updated_time":            time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &TokenString{AccessToken: "Unable to update the new data"}, err
	}
	NewaccessExpTime := time.Now().Local().Add(time.Minute * time.Duration(5)).Unix()
	log.Println("Newacceesstoken is:", NewAccessTokenString)
	log.Println(" New Token expiry time:", NewaccessExpTime)
	return &TokenString{AccessToken: NewAccessTokenString}, nil
}

type TokenString struct {
	AccessToken string `json:"accesstoken"`
}

func AutoRegenarateToken() {
	go func() {
		log.Println("Starting auto regeneration of access token...")
		for {
			time.Sleep(1 * time.Minute)
			log.Println("Refreshing the access token ..... now")
			_, err := HandleRefreshToken()
			if err != nil {
				log.Println("Auto refresh Error:", err)
			}
		}
	}()
}
