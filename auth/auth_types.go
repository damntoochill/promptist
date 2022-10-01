package auth

// VerificationResponse is what we get back
type VerificationResponse struct {
	UserID        int64
	EmailVerified bool
	Authenticated bool
}
