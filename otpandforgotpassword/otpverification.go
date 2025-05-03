package otpandforgotpassword

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func GenerateOTP() string {
	result := make([]byte, OTPlength)
	charSetLenght := big.NewInt(int64(len(otpCharSet)))
	for i := range result {
		num, _ := rand.Int(rand.Reader, charSetLenght)
		result[i] = otpCharSet[num.Int64()]
	}
	return string(result)
}

//insted of storing the otp sent to user in MongoDb or other db we can store it in Redis
//reddis will have exp time and etc...

func StoreOTPinRedis(otp string, email string, c context.Context) error {
	key := OtpKeyPrefix + email
	data, _ := bcrypt.GenerateFromPassword([]byte(otp), 10)
	res := RedisClient.Set(c, key, data, OtpExp)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func SendOTP(toName, toEmail string, otp string) error {
	fromEmail := "parshwanathparamagond1234@gmail.com"
	//fromPassword := "fbfy zhlt csqr djay"
	fromPassword := os.Getenv("KEY")
	var plainText = "Your request for Reset Password is working..."
	htmlContent := fmt.Sprintf("Hii %s \n"+
		"Your OTP for password reset is %s\r\n", toName, otp)
	var subject = "PASSWORD RESET "

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	message := gomail.NewMessage()
	message.SetHeader("From", fromEmail)
	message.SetHeader("To", toEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", plainText)
	message.AddAlternative("text/html", htmlContent)

	dialer := gomail.NewDialer(smtpHost, smtpPort, fromEmail, fromPassword)
	log.Println("Email has been sent")
	// Send Email
	return dialer.DialAndSend(message)
}







"/guest/signup": {
      "post": {
        "tags": ["guest"],
        "summary": "Guest Signup",
        "description": "Register a new guest",
        "parameters": [
          {
            "name": "guest",
            "in": "body",
            "description": "Guest Signup Details",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Guest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Guest created successfully",
            "schema": {
              "$ref": "#/definitions/GuestResponse"
            }
          },
          "400": {
            "description": "Invalid request body"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    },
    "/guest/login": {
      "post": {
        "tags": ["guest"],
        "summary": "Guest Login",
        "description": "Authenticate guest user",
        "parameters": [
          {
            "name": "guest",
            "in": "body",
            "description": "Guest login credentials",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Guest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Guest logged in successfully",
            "schema": {
              "$ref": "#/definitions/GuestResponse"
            }
          },
          "401": {
            "description": "Unauthorized - incorrect credentials"
          },
          "404": {
            "description": "Guest not found"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    },
    "/guest/logout": {
      "post": {
        "tags": ["guest"],
        "summary": "Guest Logout",
        "description": "Logs out the guest by invalidating the session token",
        "responses": {
          "200": {
            "description": "Guest logged out successfully",
            "schema": {
              "$ref": "#/definitions/GuestLogoutResponse"
            }
          },
          "400": {
            "description": "Invalid or missing token"
          },
          "401": {
            "description": "Unauthorized"
          }
        }
      }
    },
    "/guest/getallguests": {
      "get": {
        "tags": ["guest"],
        "summary": "Get All Guests",
        "description": "Retrieve all guest users from the database",
        "responses": {
          "200": {
            "description": "List of guests retrieved successfully",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Guest"
              }
            }
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    }
  },
  "definitions": {
    "Guest": {
      "type": "object",
      "required": ["email", "password", "first_name"],
      "properties": {
        "first_name": {
          "type": "string"
        },
        "last_name": {
          "type": "string"
        },
        "email": {
          "type": "string"
        },
        "password": {
          "type": "string"
        }
      }
    },
    "GuestResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        },
        "data": {
          "$ref": "#/definitions/Guest"
        }
      }
    },
    "GuestLogoutResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    }
  }