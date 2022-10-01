package post

import (
	"database/sql"
	"time"
)

// Post holds a post lawl
type Post struct {
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
	PostType        int32
	ParentID        sql.NullInt64
}

type Thread struct {
	ID              int64
	Name            string
	UserID          int64
	Username        string
	FullName        string
	Avatar          sql.NullString
	LastUserID      int64
	LastUsername    string
	LastFullName    string
	LastAvatar      sql.NullString
	NumReplies      int64
	IsSticky        bool
	IsLocked        bool
	ForumID         int64
	CreatedAt       time.Time
	CreatedAtPretty string
	UpdatedAt       time.Time
	UpdatedAtPretty string
	LastMessage     string
}

// Reply to a post or antoher reply
type Reply struct {
	ID         int64
	ParentID   int64
	UserID     int64
	Body       string
	CreatedAt  time.Time
	NumReplies int64
	NumLikes   int64
}

// Forum foo
type Forum struct {
	ID              int64
	Name            string
	Slug            string
	About           string
	NumPosts        int64
	NumViews        int64
	UpdatedAt       time.Time
	UpdatedAtPretty string
	CategoryID      int64
}

type Category struct {
	ID           int64
	Name         string
	DisplayOrder int64
}
