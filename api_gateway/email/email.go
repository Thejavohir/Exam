package email

import (
	"net/smtp"
)

func SendMail(to []string, message []byte) error {
	from := "jabdurahimov0815@gmail.com"
	password := "yqxyfopsktrbtngl"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err
}
