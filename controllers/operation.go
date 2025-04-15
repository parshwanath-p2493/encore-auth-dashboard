package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"encore.app/database"
	"encore.app/helpers"
	"encore.app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//encore:api public method=POST path=/user/signup
func Signup(ctx context.Context, user *models.Users) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	//var user models.Users
	//BodyParser(&user)
	//database.Connection()
	collection := database.OpenCollection("Users")
	id := primitive.NewObjectID()
	user.ID = id.Hex()
	user.Created_time = time.Now()
	user.Updated_time = time.Now()
	NewPassword, err := helpers.HashPsw(user.Password)
	if err != nil {
		return &Response{Message: fmt.Errorf("error in hashing the password: %w", err).Error()}, nil

	}
	user.Password = NewPassword
	err, count := helpers.Validation(user)
	if count > 0 {
		log.Println(count)
		//	log.Fatal("Enter all the required Fields", err)
		return &Response{Message: fmt.Errorf("inserting user failed: %w", err).Error()}, err
	}
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return &Response{Message: fmt.Errorf("inserting user failed: %w", err).Error()}, nil

		//return Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Error in inserting the data", "data": err.Error()})
	}
	//	return c.Status(http.StatusOK).JSON(&fiber.Map{"message": "User added successfully"})
	return &Response{Message: "User added successfully"}, nil
}

type Response struct {
	Message string `json:"message"`
}

//encore:api public method=POST path=/user/login
func Login(ctx context.Context, req *LoginReq) (*Response, error) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.Users
	collection := database.OpenCollection("Users")
	filter := bson.M{"email": req.Email}
	err := collection.FindOne(c, filter).Decode(&user)
	if err == mongo.ErrNilDocument {
		return &Response{Message: "The user does not exist"}, err
	}

	result, err := helpers.VerifyPasw(req.Password, user.Password)
	if err != nil {
		return &Response{Message: "Password is not matched"}, err
	}
	return &Response{Message: result}, nil
}

// if the error still exist the make the model for loginReqiuest
// declaring the params is very important
type LoginReq struct {
	First_Name string `json:"first_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
