package notification

import (
	"database/sql"
	"time"
)

// Notification holds a Notification lawl
type Notification struct {
	ID              int64
	UserID          int64
	FromUserID      sql.NullInt64
	FromUsername    sql.NullString
	FromFullName    sql.NullString
	FromAvatar      sql.NullString
	CreatedAt       time.Time
	CreatedAtPretty string
	NotifType       int32
	ContentID       sql.NullInt64
	Body            sql.NullString
	ArtUUID         sql.NullString
}
