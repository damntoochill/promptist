package royale

import (
	"database/sql"
	"time"
)

type Royale struct {
	ID                int64
	Name              string
	Description       string
	Slug              string
	CreatedAt         time.Time
	CreatedAtPretty   string
	NumSubmissions    int64
	Winner            sql.NullInt64
	WinningSubmission sql.NullString
	CanSubmit         bool
	CanVote           bool
	Prize             string
	Round1Body        string
	Round2Body        string
	Round3Body        string
	Round4Body        string
}

type Submission struct {
	ID              int64
	RoyaleID        int64
	UserID          int64
	ArtID           string
	CreatedAt       time.Time
	CreatedAtPretty string
}

type Vote struct {
	ID           int64
	RoyaleID     int64
	UserID       int64
	SubmissionID int64
}
