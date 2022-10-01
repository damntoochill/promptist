package usersession

// UserSession stores data that we need consistently
type UserSession struct {
	UserID              int64
	Username            string
	IsAuthenticated     bool
	IsVerified          bool
	IsAdmin             bool
	UnreadNotifications bool
	UnreadChats         bool
}
