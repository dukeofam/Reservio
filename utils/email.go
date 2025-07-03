package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

var (
	SMTPHost     = "smtp.gmail.com"
	SMTPPort     = "587"
	SMTPUser     = "your_email@gmail.com" // Change for production
	SMTPPassword = "your_password"        // Change for production
)

// NOTE: In production, SendMail should be mocked in tests to avoid sending real emails.
func SendMail(to, subject, body string) error {
	// Skip real SMTP in test mode
	if os.Getenv("TEST_MODE") == "1" {
		log.Printf("[SendMail] Test mode - skipping real email to %s: %s", to, subject)
		return nil
	}

	from := SMTPUser
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	auth := smtp.PlainAuth("", SMTPUser, SMTPPassword, SMTPHost)
	addr := fmt.Sprintf("%s:%s", SMTPHost, SMTPPort)
	var err error
	for i := 1; i <= 3; i++ {
		err = smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
		if err == nil {
			return nil
		}
		log.Printf("[SendMail] Attempt %d failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	log.Printf("[SendMail] All attempts failed: %v", err)
	return err
}
