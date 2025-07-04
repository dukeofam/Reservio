package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

var (
	smtpHostEnv = "SMTP_HOST"
	smtpPortEnv = "SMTP_PORT"
	smtpUserEnv = "SMTP_USER"
	smtpPassEnv = "SMTP_PASSWORD"
)

// NOTE: In production, SendMail should be mocked in tests to avoid sending real emails.
func SendMail(to, subject, body string) error {
	// Skip real SMTP in test mode but return a dummy error so callers can assert failure without sending mail
	if os.Getenv("TEST_MODE") == "1" {
		log.Printf("[SendMail] Test mode - skipping real email to %s: %s", to, subject)
		return fmt.Errorf("send mail skipped in test mode")
	}

	// Load SMTP config from environment (with sane defaults for local dev)
	host := getenvDefault(smtpHostEnv, "smtp.gmail.com")
	port := getenvDefault(smtpPortEnv, "587")
	user := os.Getenv(smtpUserEnv)
	pass := os.Getenv(smtpPassEnv)

	// In production, require explicit credentials
	if os.Getenv("ENVIRONMENT") == "production" {
		if user == "" || pass == "" {
			return fmt.Errorf("SMTP credentials not set (SMTP_USER/PASSWORD)")
		}
	}

	// Fallback defaults for non-production
	if user == "" {
		user = "your_email@gmail.com"
	}

	from := user
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)
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

// getenvDefault returns env value or default if unset
func getenvDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
