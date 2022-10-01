package message

import (
	"database/sql"
	"time"

	"github.com/promptist/web/profile"
)

// Message holds a Message lawl
type Message struct {
	ID              int64
	UserID          int64
	Username        string
	FullName        string
	UserBadges      sql.NullString
	Avatar          sql.NullString
	Body            string
	CreatedAt       time.Time
	CreatedAtPretty string
	ChatID          int64
}

// Chat holds all the messages for a chat
type Chat struct {
	ID              int64
	CreatedAt       time.Time
	CreatedAtPretty string
	UpdatedAt       time.Time
	UpdatedAtPretty string
	LastMessage     sql.NullString
	IsRead          bool
	Members         map[int64]profile.Brief
	IsGroup         bool
	Recipient       profile.Brief
}
