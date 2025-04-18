package helpers

import (
	"errors"
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

func GenerateJwt(name string, email string) (string, string, error) {
	claims := &Info{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(24)).Unix(), //Create the access token (expires in 24 minutes)
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

	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":   name,
		"expire": jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()},
	})
	SignedRefreshToken, err := RefreshToken.SignedString([]byte(Secret_key))
	if err != nil {
		return "There is error in signing the token", "", err
	}
	log.Println("Signed token accesstoken is :", SignedAccessToken)

	return SignedAccessToken, SignedRefreshToken, nil
}
func RefreshToken(refreshTokenString string) (string, error) {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fetch the secret key from .env file
	SecretKey := os.Getenv("SECRET_KEY")

	// Parse and verify the refresh token
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(SecretKey), nil
	})

	// Check if the refresh token is valid
	if err != nil || !refreshToken.Valid {
		return "", errors.New("invalid refresh token")
	}

	// Extract claims from the refresh token
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		// Extract the name from the claims
		name, ok := claims["name"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}

		// Create a new access token with a 15-minute expiration time
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name": name,
			"exp":  time.Now().Add(time.Minute * 15).Unix(), // New access token expires in 15 minutes
		})

		// Sign and generate the new access token string
		newAccessTokenString, err := newAccessToken.SignedString([]byte(SecretKey))
		if err != nil {
			return "", errors.New("failed to generate new access token")
		}

		// Return the new access token string
		return newAccessTokenString, nil
	}

	// If we couldn't extract claims or the token is invalid
	return "", errors.New("invalid refresh token")
}
