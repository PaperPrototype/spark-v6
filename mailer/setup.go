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
}

func SendgridTestEmail() {
	from := mail.NewEmail("Example User", "info@sparker3d.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example Sparker", "a.spytech360@gmail.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(helpers.GetSendgridKey())
	response, err := client.Send(message)
	if err != nil {
		log.Println("failed")
		log.Println(err)
	} else {
		log.Println("even if not verified identity at least we didn't get an error. So => SUCCESS.")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
