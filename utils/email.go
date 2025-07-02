package utils

import (
	"fmt"
	"net/smtp"
)

var (
	SMTPHost     = "smtp.gmail.com"
	SMTPPort     = "587"
	SMTPUser     = "your_email@gmail.com" // Change for production
	SMTPPassword = "your_password"        // Change for production
)

// NOTE: In production, SendMail should be mocked in tests to avoid sending real emails.
func SendMail(to, subject, body string) error {
	from := SMTPUser
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	auth := smtp.PlainAuth("", SMTPUser, SMTPPassword, SMTPHost)
	addr := fmt.Sprintf("%s:%s", SMTPHost, SMTPPort)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
