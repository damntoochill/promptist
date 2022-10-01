package email

import (
	"context"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var ec = Config{
	Domain:       os.Getenv("MAILGUN_DOMAIN"),
	APIKey:       os.Getenv("MAILGUN_API_KEY"),
	PublicAPIKey: os.Getenv("MAILGUN_PUBLIC_API_KEY"),
	EmailSender:  os.Getenv("EMAIL_SENDER"),
	AppName:      os.Getenv("APP_NAME"),
	WebsiteHost:  os.Getenv("WEBSITE_HOST"),
}

// Email a user through Mailgun API
func Email(subject string, body string, recipient string) error {
	mg := mailgun.NewMailgun(ec.Domain, ec.APIKey)
	message := mg.NewMessage(ec.EmailSender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

// PasswordReset sends the password reset email
func PasswordReset(email string, token string) error {
	subject := "Password reset"
	body := "A password reset has been initiated on your acount. Click the link below to reset your password. If you didn't request one, somebody else did! RESET LINK: " + ec.WebsiteHost + "/reset-password/?token=" + token + "&email=" + email
	err := Email(subject, body, email)
	if err != nil {
		return err
	}
	return nil
}

// Verification emails a token to verify the account
func Verification(email string, token string) error {
	subject := "Welcome to " + ec.AppName + "!"
	body := "Hello from " + ec.AppName + "! Verify your email by clicking the following link: " + ec.WebsiteHost + "/email-verification/?token=" + token + "&email=" + email
	err := Email(subject, body, email)
	if err != nil {
		return err
	}
	return nil
}

// ResendVerification resens a verification token
func ResendVerification(email string, token string) error {
	subject := "Verify your account on " + ec.AppName + "!"
	body := "Hello from " + ec.AppName + "! Verify your email by clicking the following link: " + ec.WebsiteHost + "/email-verification/?token=" + token + "&email=" + email
	err := Email(subject, body, email)
	if err != nil {
		return err
	}
	return nil
}
