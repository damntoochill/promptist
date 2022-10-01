package comment

import (
	"database/sql"
	"os"

	"github.com/promptist/web/art"
	"github.com/promptist/web/auth"
	"github.com/promptist/web/email"
	"github.com/promptist/web/notification"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/promptist/web/utils"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// NewComment creates a new comment
func NewComment(db *sql.DB, body string, piece art.Piece, us *usersession.UserSession) (int64, bool) {
	commenterBrief, _ := profile.GetBrief(db, us.UserID)

	stmt, err := db.Prepare("INSERT INTO comment_comments (body, user_id, username, full_name, avatar, piece_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(body, us.UserID, commenterBrief.Username, commenterBrief.FullName, commenterBrief.PhotoUUID, piece.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	commentID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	updateCommentCount(db, piece.ID)
	b := utils.NewNullString(body)
	artUUID := utils.NewNullString(piece.ImageUUID)
	ok := notification.NewNotification(db, b, sql.NullInt64{}, artUUID, 0, commenterBrief, piece.UserID)
	if !ok {
		log.Error().Msg("unable to create notification")
	}

	ownerEmail, err := auth.Email(db, piece.UserID)
	if err == nil {
		var subject string
		if piece.Name.Valid {
			subject = "New comment on " + piece.Name.String
		} else {
			subject = "New comment on your art!"
		}
		emailBody := "Facts. You've got a new comment on one of your pieces! You can view it here: " + os.Getenv("WEBSITE_HOST") + "/art/" + piece.ImageUUID

		err = email.Email(subject, emailBody, ownerEmail)
		if err != nil {
			log.Error().Err(err).Msg("unable to send email")
		}
	}

	return commentID, true
}

// GetComments from the database
func GetComments(db *sql.DB, pieceID int64) ([]Comment, bool) {
	q := `
		SELECT comment_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, parent_id, piece_id
		FROM comment_comments
		WHERE piece_id = ?
		ORDER BY created_at
		`
	res, err := db.Query(q, pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Comment{}, false
	}
	defer res.Close()

	var comments []Comment

	for res.Next() {
		var c Comment
		err := res.Scan(&c.ID, &c.Body, &c.CreatedAt, &c.UserID, &c.Username, &c.FullName, &c.Avatar, &c.NumReplies, &c.NumLikes, &c.ParentID, &c.PieceID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Comment{}, false
		}
		c.CreatedAtPretty = timeago.New(c.CreatedAt).Format()
		comments = append(comments, c)
	}
	return comments, true
}

func updateCommentCount(db *sql.DB, pieceID int64) bool {
	log.Debug().Msg("UpdateCommentCount")
	stmt, err := db.Prepare("UPDATE art_pieces SET comments = comments + 1 WHERE piece_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}
