package comment

import (
	"database/sql"
	"os"

	"github.com/promptist/web/auth"
	"github.com/promptist/web/email"
	"github.com/promptist/web/notification"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/promptist/web/utils"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// NewProfileComment creates a new comment
func NewProfileComment(db *sql.DB, body string, profileUserID int64, us *usersession.UserSession) (int64, bool) {
	commenterBrief, _ := profile.GetBrief(db, us.UserID)

	stmt, err := db.Prepare("INSERT INTO comment_pro_comments (body, user_id, username, full_name, avatar, profile_user_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(body, us.UserID, commenterBrief.Username, commenterBrief.FullName, commenterBrief.PhotoUUID, profileUserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	commentID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	updateProfileCommentCount(db, profileUserID)

	b := utils.NewNullString(body)
	ok := notification.NewNotification(db, b, sql.NullInt64{}, sql.NullString{}, 1, commenterBrief, profileUserID)
	if !ok {
		log.Error().Msg("unable to create notification")
	}

	ownerEmail, err := auth.Email(db, profileUserID)
	ownerBrief, _ := profile.GetBrief(db, profileUserID)
	if err == nil {
		subject := "New profile comment from " + commenterBrief.FullName

		emailBody := commenterBrief.FullName + " left a comment on your profile! You can view it here: " + os.Getenv("WEBSITE_HOST") + "/people/" + ownerBrief.Username

		err = email.Email(subject, emailBody, ownerEmail)
		if err != nil {
			log.Error().Err(err).Msg("unable to send email")
		}
	}

	return commentID, true
}

// GetProfileComments from the database
func GetProfileComments(db *sql.DB, profileUserID int64) ([]ProfileComment, bool) {
	q := `
		SELECT pro_comment_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, parent_id, profile_user_id
		FROM comment_pro_comments
		WHERE profile_user_id = ?
		ORDER BY created_at DESC
		`
	res, err := db.Query(q, profileUserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []ProfileComment{}, false
	}
	defer res.Close()

	var comments []ProfileComment

	for res.Next() {
		var c ProfileComment
		err := res.Scan(&c.ID, &c.Body, &c.CreatedAt, &c.UserID, &c.Username, &c.FullName, &c.Avatar, &c.NumReplies, &c.NumLikes, &c.ParentID, &c.ProfileUserID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []ProfileComment{}, false
		}
		c.CreatedAtPretty = timeago.New(c.CreatedAt).Format()
		comments = append(comments, c)
	}
	return comments, true
}

func updateProfileCommentCount(db *sql.DB, profileUserID int64) bool {
	log.Debug().Msg("UpdateCommentCount")
	stmt, err := db.Prepare("UPDATE profile_profiles SET num_comments = num_comments + 1 WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(profileUserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}
