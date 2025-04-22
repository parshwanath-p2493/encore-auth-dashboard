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
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fromName = os.Getenv("FROM_NAME")
	fromEmail = os.Getenv("FROM_EMAIL")
}

func SendMail(toName, toEmail string) (string, error) {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	subject := "Sending with Twilio SendGrid is Fun"
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("SendGrid error:", err)
		return "", err
	}

	fmt.Println("Status:", response.StatusCode)
	fmt.Println("Body:", response.Body)
	fmt.Println("Headers:", response.Headers)

	return "Welcome message sent successfully", nil
}
