package email

import (
	"net/smtp"
)

func SendMail(to []string, message []byte) error {
	from := "javohirabdurahimovcom@gmail.com"
	password := "ihtizrusjypnultt"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err
}