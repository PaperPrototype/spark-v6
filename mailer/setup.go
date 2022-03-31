package mailer

import (
	"fmt"
	"log"
	"main/helpers"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var client *sendgrid.Client

func Setup() {
	client = sendgrid.NewSendClient(helpers.GetSendgridKey())

	SendgridTestEmail()
}

func SendgridTestEmail() {
	from := mail.NewEmail("Sparker", "spark3dsoftware@gmail.com")

	to := mail.NewEmail("You", "a.spytech360@gmail.com")

	subject := "Sending with SendGrid is Fun"
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("success!")
		fmt.Println(response.StatusCode)
	}
}
