package email

// Config holds the things we need to send that
// ditigal goodness through the MailGun API
type Config struct {
	Domain       string
	APIKey       string
	PublicAPIKey string
	EmailSender  string
	AppName      string
	WebsiteHost  string
}
