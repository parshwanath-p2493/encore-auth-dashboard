package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	fromName  string
	fromEmail string
	key       string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fromName = os.Getenv("FROM_NAME")
	fromEmail = os.Getenv("FROM_EMAIL")
	key = os.Getenv("SENDGRID_API_KEY")

}

func SendMail(toName, toEmail string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	subject := "Welcome Message"
	plainTextContent := "and easy to do anywhere, even with Go"
	//htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	var htmlContent = fmt.Sprint("<strong> 🔔 HI WELCOME TO OUR APP</strong>")

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(key)
	response, err := client.Send(message)
	if err != nil {
		log.Println("SendGrid error:", err)
		return err
	}

	fmt.Println("Status:", response.StatusCode)
	fmt.Println("Body:", response.Body)
	fmt.Println("Headers:", response.Headers)
	log.Println("Welcome message sent successfully")
	return nil
}

/**
var plainText = "and easy to do anywhere, even with Go"
var htmlContent = fmt.Sprint("<strong> 🔔 HI WELCOME TO OUR APP</strong>")

var subject = "WEKCOME MESSAGE "

func SendMail(toName, toEmail string) error {
	fromEmail := "parshwanathparamagond1234@gmail.com"
	//fromPassword := "fbfy zhlt csqr djay"
	fromPassword := os.Getenv("KEY")

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
**/
