package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
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
	var ExistingUser models.Users
	filter := bson.M{"email": user.Email}
	err := collection.FindOne(ctx, filter).Decode(&ExistingUser)
	if err == nil {
		return &Response{Message: "The user is already existed...."}, err
	}
	if err != mongo.ErrNoDocuments {
		return &Response{Message: "Error while checking existing user"}, err
	}
	id := primitive.NewObjectID()
	user.ID = id.Hex()
	user.Created_time = time.Now().Local()
	user.Updated_time = time.Now().Local()
	_, err = helpers.EmailValidation(user.Email)
	if err != nil {
		return &Response{Message: "INVALID EMAIL"}, err
	}
	NewPassword, err := helpers.HashPsw(user.Password)
	if err != nil {
		return &Response{Message: fmt.Errorf("error in hashing the password: %w", err).Error()}, nil

	}
	user.Password = NewPassword
	Token, Refresh_Token, err := helpers.GenerateJwt(user.First_Name, user.Email)
	if err != nil {
		return &Response{Message: fmt.Errorf("error in creating the new token: %w", err).Error()}, nil
	}
	user.Token = Token
	user.Refresh_Token = Refresh_Token
	// helpers.GlobalSignedAccessToken = Token
	// helpers.GlobalSignedRefreshToken = Refresh_Token
	//helpers.HandleRefreshToken()
	go helpers.AutoRegenarateToken()

	//user.Email = user.Email
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
	err = helpers.SendMail(user.First_Name, user.Email)
	if err != nil {
		return &Response{Message: "There error is updating "}, err
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
	_, err := helpers.EmailValidation(req.Email)
	if err != nil {
		return &Response{Message: "The mail is not valid enter valid.."}, err
	}
	filter := bson.M{"email": req.Email}
	err = collection.FindOne(c, filter).Decode(&user)
	if err == mongo.ErrNilDocument {
		return &Response{Message: "The user does not exist"}, err
	}

	result, err := helpers.VerifyPasw(req.Password, user.Password)
	if err != nil {
		return &Response{Message: "Password is not matched"}, err
	}
	accesstoken, refreshtoken, err := helpers.GenerateJwt(user.First_Name, user.Email)
	if err != nil {
		return &Response{Message: "There is error in login try again..."}, err
	}
	//helpers.HandleRefreshToken()
	go helpers.AutoRegenarateToken()
	update := bson.M{
		"$set": bson.M{
			"token":         accesstoken,
			"refresh_token": refreshtoken,
			"updated_time":  time.Now().Local(),
		},
	}
	_, err = collection.UpdateOne(c, filter, update)
	if err != nil {
		return &Response{Message: "There error is updating "}, err
	}

	err = helpers.SendMail(user.First_Name, user.Email)
	if err != nil {
		return &Response{Message: "There error is updating "}, err
	}
	return &Response{Message: result}, nil
}

// if the error still exist the make the model for loginReqiuest
// declaring the params is very important
type LoginReq struct {
	First_Name string `json:"first_name"` //Always use Uppercase for exporting the variables.../
	Email      string `json:"email"`
	Password   string `json:"password"`
} //API request types must be named structs.        //Need to note important

//encore:api auth method=POST path=/user/delete
func DeleteUser(ctx context.Context, req *DeleteUserReq) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	inputEmail := strings.TrimSpace(strings.ToLower(req.Email))
	log.Printf("Received request to delete user. Normalized email: '%s'", inputEmail)

	if inputEmail == "" {
		return &Response{Message: "Invalid email provided"}, nil
	}

	collection := database.OpenCollection("Users")
	filter := bson.M{"email": inputEmail}

	var user models.Users
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println("User not found before deletion attempt:", err)
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return &Response{Message: "Error in deleting the user"}, err
	}
	if result.DeletedCount == 0 {
		return &Response{Message: "User not found, unable to delete"}, nil
	}

	return &Response{Message: "User successfully deleted"}, nil
}

type DeleteUserReq struct {
	UserName string `json:"username"`
	Email    string `json:"email"` // preferred
}

//encore:api public method=POST path=/user/refreshtoken
func RefreshToken(ctx context.Context, req *RequestRefresh) (*AccessTokenResponse, error) {
	newToken, err := helpers.HandleRefreshToken()
	if err != nil {
		return &AccessTokenResponse{Message: "Unable to generate accesss token or token not yet expired", AccessToken: ""}, err
	}
	return &AccessTokenResponse{Message: ".............d ", AccessToken: newToken.AccessToken}, nil

}

type RequestRefresh struct {
	RefreshToken string `json:"refreshtoken"`
}
type AccessTokenResponse struct {
	Message     string
	AccessToken string
}

type InputForgotPassword struct {
	Email string `json:"email"`
}

//encore:api public method=POST path=/user/forgotpassword
func ForgotPassword(ctx context.Context, input *InputForgotPassword) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return &Response{Message: "The Otp is sent to your email..	"}, nil
}
