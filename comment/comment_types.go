package comment

import (
	"database/sql"
	"time"
)

// Comment holds a Comment lawl
type Comment struct {
	ID              int64
	UserID          int64
	Username        string
	FullName        string
	UserBadges      sql.NullString
	Avatar          sql.NullString
	Body            string
	CreatedAt       time.Time
	CreatedAtPretty string
	NumReplies      int64
	NumLikes        int64
	ParentID        sql.NullInt64
	PieceID         int64
}

// ProfileComment holds a ProfileComment lawl
type ProfileComment struct {
	ID              int64
	UserID          int64
	Username        string
	FullName        string
	UserBadges      sql.NullString
	Avatar          sql.NullString
	Body            string
	CreatedAt       time.Time
	CreatedAtPretty string
	NumReplies      int64
	NumLikes        int64
	ParentID        sql.NullInt64
	ProfileUserID   int64
}
