package helpers

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"encore.app/database"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var (
// GlobalSignedAccessToken  string // Removed: Tokens should not be global
// GlobalSignedRefreshToken string // Removed: Tokens should not be global
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
	log.Println("Access token expire time is :", claims.ExpiresAt)
	//we created the refresh token here ... expire time is 24 hours now..

	Refreshclaims := &Info{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(), //Create the access token (expires in 5 minutes)
		},
	}
	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Refreshclaims)
	SignedRefreshToken, err := RefreshToken.SignedString([]byte(Secret_key))
	if err != nil {
		return "There is error in signing the token", "", err
	}
	log.Println("Signed token refreshtoken is :", SignedRefreshToken)
	log.Println("Refresh token expire time is :", Refreshclaims.ExpiresAt)
	// GlobalSignedAccessToken = SignedAccessToken // Removed
	// GlobalSignedRefreshToken = SignedRefreshToken // Removed

	return SignedAccessToken, SignedRefreshToken, nil
}

// HandleRefreshToken validates a given refresh token and issues a new access token.
// It's intended to be called by an API endpoint when a client requests token refresh.
func HandleRefreshToken(providedRefreshTokenString string) (*TokenString, error) {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file") // Consider returning error instead of fatal
	}
	log.Println("Attempting to refresh token")

	// Fetch the secret key from .env file
	SecretKey := os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		log.Println("SECRET_KEY not found in environment variables")
		return nil, errors.New("server configuration error: missing secret key")
	}

	Refreshclaims := &Info{}
	RefreshToken, err := jwt.ParseWithClaims(providedRefreshTokenString, Refreshclaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		log.Println("Error parsing refresh token:", err)
		return &TokenString{AccessToken: "Invalid refresh token"}, err // More specific error might be useful
	}

	if !RefreshToken.Valid {
		log.Println("Refresh token is invalid or expired")
		return &TokenString{AccessToken: "Refresh token invalid or expired"}, errors.New("refresh token invalid or expired")
	}

	log.Println("Refresh token validated successfully.")
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
	// GlobalSignedAccessToken = NewAccessTokenString // Removed
	if err != nil {
		log.Print("Failed to sign new access token:", err)
		return nil, err
	}

	OriginalRefreshExpireTime := time.Unix(Refreshclaims.ExpiresAt, 0) // Expiry of the provided refresh token
	log.Println("Original Refresh Token expiry time", OriginalRefreshExpireTime)
	NewAccessExpireTime := time.Unix(NewClaims.ExpiresAt, 0)
	log.Println("New Access Token generated, expiry time:", NewAccessExpireTime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.OpenCollection("Users")
	filter := bson.M{"email": Refreshclaims.Email} // Ensure this email is correct and from validated token
	update := bson.M{
		"$set": bson.M{
			"token":           NewAccessTokenString,
			"expiretimetoken": NewAccessExpireTime,
			// "expiretimerefresh_token": OriginalRefreshExpireTime, // Keep original refresh token's expiry
			// "refresh_token":           providedRefreshTokenString, // Store the same refresh token, assuming no rotation
			"updated_time": time.Now().Local(),
		},
		// Optionally, if you want to update the refresh token details (e.g. if it was rotated)
		// "$set": bson.M{
		// 	"token": NewAccessTokenString,
		// 	"expiretimetoken": NewAccessExpireTime,
		// 	"refresh_token": newRefreshTokenString, // if you implement refresh token rotation
		// 	"expiretimerefresh_token": newRefreshExpireTime, // if you implement refresh token rotation
		// 	"updated_time": time.Now().Local(),
		// },
	}
	// Only update the access token related fields, leave refresh token as is unless rotating
	// If you are not rotating refresh tokens, you might not even need to update the refresh_token field here
	// if it's already stored correctly. The main thing is to update the access token.
	// For simplicity, the current update only sets the new access token and its expiry.
	// You might want to ensure `refresh_token` and `expiretimerefresh_token` are correctly set
	// in the database upon initial login.

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

// func AutoRegenarateToken() { // Removed: This approach is not suitable for user-specific tokens
// 	go func() {
// 		log.Println("Starting auto regeneration of access token...")
// 		for {
// 			time.Sleep(1 * time.Minute)
// 			log.Println("Refreshing the access token ..... now")
// 			// _, err := HandleRefreshToken() // This would need a refresh token passed in
// 			// if err != nil {
// 			// 	log.Println("Auto refresh Error:", err)
// 			// }
// 		}
// 	}()
// }
