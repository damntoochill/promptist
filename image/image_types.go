package image

// Config holds options that we need to connect to our 3rd pary
// storage solution
type Config struct {
	Key      string
	Secret   string
	Endpoint string
	Region   string
	Bucket   string
}
