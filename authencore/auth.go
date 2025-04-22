package authencore

import (
	"context"
	"errors"
	"log"
	"os"

	"encore.dev/beta/auth"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type User struct {
	Name  string
	Email string
}
type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, *User, error) {
	var err error = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load the environmental values...")
	}
	Secret_key := os.Getenv("SECRET_KEY")
	claims := &Claims{}
	result, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(Secret_key), nil
	})
	if err != nil || !result.Valid {
		return "", nil, errors.New("invalid Token or expired token")
	}
	return auth.UID(claims.Email), &User{
		Name:  claims.Name,
		Email: claims.Email,
	}, nil
}
