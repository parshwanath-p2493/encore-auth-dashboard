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
