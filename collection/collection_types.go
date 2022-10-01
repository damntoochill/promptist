package collection

import (
	"database/sql"
	"time"
)

// Collection holds a collection of art
type Collection struct {
	ID          int64
	Name        string
	Description sql.NullString
	CreatedAt   time.Time
	IsPublic    bool
	UserID      int64
	NumPieces   int64
	PieceID     int64
}
