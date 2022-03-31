package mailer

import (
	"fmt"
	"log"
	"main/db"
	"main/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// in minutes
const MinUntilVerifyExpires time.Duration = 10

// generate verification link and send email
func SendVerification(userID uint64) error {
	user, err1 := db.GetUser(userID)
	if err1 != nil {
		log.Println("mailer ERROR getting user:", err1)
		return err1
	}

	verify := db.Verify{
		UserID:     user.ID,
		VerifyUUID: uuid.NewString(),
		ExpiresAt:  time.Now().Add(time.Minute * MinUntilVerifyExpires),
	}
	err3 := db.CreateVerify(&verify)
	if err3 != nil {
		log.Println("mailer ERROR creating verification link:", err3)
		return err3
	}
	htmlContent := buildEmail(
		"Verify your account",
		"This email and link is to verify your new account with Sparker.",
		helpers.GetHost()+"/login/verify/"+verify.VerifyUUID,
		"If you did not sign up for an account on Sparker.com you can safely ignore this email.",
	)

	from := mail.NewEmail("Sparker", "spark3dsoftware@gmail.com")

	to := mail.NewEmail("You", "a.spytech360@gmail.com")

	subject := "Sending with SendGrid is Fun"
	plainTextContent := "and easy to do anywhere, even with Go"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	response, err := client.Send(message)
	if err != nil {
		log.Println("mailer/emails ERROR sending verification email", err)
		return err
	} else {
		log.Println("mailer/emails success!")
		fmt.Println(response.StatusCode)
	}

	return nil
}

func buildEmail(title string, message string, link string, afterLinkMessage string) string {
	return `
	<html>
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		</head>
	<body>
		<div style="padding-right: 5%; padding-left: 5%;">
			<h2>` + title + `</h2>
			<p>
				` + message + `
			</p>
			<p>
				<a href="` + link + `">` + link + `</a>
			</p>
			<p>
				` + afterLinkMessage + `
			</p>
		</div>
	</body>`
}
