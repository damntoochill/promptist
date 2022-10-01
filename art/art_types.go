package art

import (
	"database/sql"
	"time"
)

// Piece is the omega
type Piece struct {
	ID                    int64
	Name                  sql.NullString
	Slug                  sql.NullString
	Description           sql.NullString
	ImageUUID             string
	Prompt                sql.NullString
	IsDraft               bool
	Likes                 int32
	Views                 int32
	Saves                 int32
	Comments              int32
	CreatedAt             time.Time
	Tags                  []Tag
	TagsLiteral           sql.NullString
	UserID                int64
	Username              string
	FullName              string
	UserBadges            sql.NullString
	ProfilePhotoUUID      sql.NullString
	ProgramID             int64
	ProgramName           sql.NullString
	ProgramSlug           sql.NullString
	ProgramCoverImageUUID sql.NullString
}

// Tag is a way to navigate this crazy world
type Tag struct {
	ID          int64
	Name        string
	Total       int32
	Description sql.NullString
}

// Program is the name of the AI that created the art
type Program struct {
	ID             int64
	Name           string
	Slug           string
	Description    string
	CoverImageUUID string
	IsOpenSource   bool
}

// Like foo
type Like struct {
	ID        int64
	PieceID   int64
	UserID    int64
	CreatedAt time.Time
}
