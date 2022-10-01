package notification

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// GetNotifTypeString returns the type in string format
// This is more for debugging
func GetNotifTypeString(notifType int32) string {
	switch notifType {
	case 0:
		return "art comment"
	case 1:
		return "profile comment"
	case 2:
		return "follow"
	case 3:
		return "like"
	}
	return "error"
}

// NewNotification creates a new post
func NewNotification(db *sql.DB, body sql.NullString, contentID sql.NullInt64, artUUID sql.NullString, notifType int64, fromBrief profile.Brief, toUserID int64) bool {

	stmt, err := db.Prepare("INSERT INTO notif_notifications (user_id, from_user_id, from_username, from_full_name, from_avatar, notif_type, content_id, art_uuid, body) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(toUserID, fromBrief.UserID, fromBrief.Username, fromBrief.FullName, fromBrief.PhotoUUID, notifType, contentID, artUUID, body)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	updateNotificationCount(db, toUserID)

	return true
}

// GetNotifications from the database
func GetNotifications(db *sql.DB, userID int64) ([]Notification, bool) {
	q := `
		SELECT id, user_id, from_user_id, from_username, from_full_name, from_avatar, created_at, notif_type, content_id, art_uuid, body
		FROM notif_notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Notification{}, false
	}
	defer res.Close()

	var notifications []Notification

	for res.Next() {
		var n Notification
		err := res.Scan(&n.ID, &n.UserID, &n.FromUserID, &n.FromUsername, &n.FromFullName, &n.FromAvatar, &n.CreatedAt, &n.NotifType, &n.ContentID, &n.ArtUUID, &n.Body)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Notification{}, false
		}
		n.CreatedAtPretty = timeago.New(n.CreatedAt).Format()
		notifications = append(notifications, n)
	}
	return notifications, true
}

// NumNotifications returns the post from the database
func NumNotifications(db *sql.DB, userID int64) (int64, bool) {
	stmt, err := db.Prepare("SELECT num_notifications FROM profile_profiles WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return 0, false
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, false
	}
	var n int64
	err = rows.Scan(&n)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return 0, false
	}
	return n, true
}

func updateNotificationCount(db *sql.DB, userID int64) bool {
	log.Debug().Msg("UpdateNotificationCount")
	stmt, err := db.Prepare("UPDATE profile_profiles SET num_notifications = num_notifications + 1, unread_notifications=1 WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}

// UnreadNotifications returns the post from the database
func UnreadNotifications(db *sql.DB, userID int64) (bool, bool) {
	stmt, err := db.Prepare("SELECT unread_notifications FROM profile_profiles WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false, false
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false, false
	}
	defer rows.Close()
	if !rows.Next() {
		return false, false
	}
	var unread bool
	err = rows.Scan(&unread)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return false, false
	}
	return unread, true
}

func SetNotificationsToRead(db *sql.DB, userID int64) bool {
	log.Debug().Msg("SetNotificationsToRead")
	stmt, err := db.Prepare("UPDATE profile_profiles SET unread_notifications=0 WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}

func CheckForUnread(c *gin.Context, db *sql.DB, us *usersession.UserSession) bool {
	log.Debug().Msg("CheckForUnread")

	unread, ok := UnreadNotifications(db, us.UserID)
	if !ok {
		return false
	}

	usersession.Save(c, us.IsAdmin, us.IsAuthenticated, us.IsVerified, us.UserID, us.Username, unread, us.UnreadNotifications)

	return true
}
