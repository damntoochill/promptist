package main

import (
	"errors"
	"os"
)

func checkConfig() error {
	if os.Getenv("PORT") == "" {
		return errors.New("$PORT must be set")
	} else if os.Getenv("APP_NAME") == "" {
		return errors.New("$APP_NAME must be set")
	} else if os.Getenv("DB_DATABASE") == "" {
		return errors.New("$DB_DATABASE must be set")
	} else if os.Getenv("DB_USER") == "" {
		return errors.New("$DB_USER must be set")
	} else if os.Getenv("DB_PASSWORD") == "" {
		return errors.New("$DB_PASSWORD must be set")
	} else if os.Getenv("DB_HOST") == "" {
		return errors.New("$DB_HOST must be set")
	} else if os.Getenv("DB_PORT") == "" {
		return errors.New("$DB_PORT must be set")
	} else if os.Getenv("SESSION_KEY") == "" {
		return errors.New("$SESSION_KEY must be set")
	} else if os.Getenv("MAILGUN_DOMAIN") == "" {
		return errors.New("$MAILGUN_DOMAIN must be set")
	} else if os.Getenv("MAILGUN_API_KEY") == "" {
		return errors.New("$MAILGUN_API_KEY must be set")
	} else if os.Getenv("MAILGUN_PUBLIC_API_KEY") == "" {
		return errors.New("$MAILGUN_PUBLIC_API_KEY must be set")
	} else if os.Getenv("EMAIL_SENDER") == "" {
		return errors.New("$EMAIL_SENDER must be set")
	} else if os.Getenv("WEBSITE_HOST") == "" {
		return errors.New("$WEBSITE_HOST must be set")
	}
	return nil
}
