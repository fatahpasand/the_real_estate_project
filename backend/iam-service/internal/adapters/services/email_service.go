// internal/adapters/services/email_service.go
package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type emailService struct {
	from     string
	password string
	host     string
	port     string
}

func NewEmailService() *emailService {
	return &emailService{
		from:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		host:     "smtp-mail.outlook.com",
		port:     "587",
	}
}

func (s *emailService) SendVerificationEmail(to, otp string) error {
	subject := "Email Verification"
	body := fmt.Sprintf("Your verification code is: %s\nValid for 15 minutes.", otp)
	msg := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func (s *emailService) SendLoginAlert(to, ip, userAgent string) error {
	subject := "New Login Detected"
	body := fmt.Sprintf("New login detected from:\nIP: %s\nDevice: %s", ip, userAgent)
	msg := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
