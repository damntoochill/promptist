package post

import (
	"database/sql"

	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// NewReply creates a new reply to a post
func NewReply(db *sql.DB, body string, post int64, us *usersession.UserSession) (int64, bool) {
	log.Debug().Msg("NewReply")
	stmt, err := db.Prepare("INSERT INTO post_posts (body, parent_id, user_id, username, full_name, avatar, post_type) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	brief, _ := profile.GetBrief(db, us.UserID)
	res, err := stmt.Exec(body, post, us.UserID, brief.Username, brief.FullName, brief.PhotoUUID, 1)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	replyID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get insert ID")
		return 0, false
	}
	if !UpdateReplyCount(db, post) {
		return 0, false
	}
	return replyID, true
}

func GetReplies(db *sql.DB, post int64) ([]Post, bool) {
	log.Debug().Msg("GetReplies")
	stmt, err := db.Prepare("SELECT post_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, post_type, parent_ID FROM post_posts WHERE parent_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return nil, false
	}
	rows, err := stmt.Query(post)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return nil, false
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.Body, &p.CreatedAt, &p.UserID, &p.Username, &p.FullName, &p.Avatar, &p.NumReplies, &p.NumLikes, &p.PostType, &p.ParentID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan row")
			return nil, false
		}
		p.CreatedAtPretty = timeago.New(p.CreatedAt).Format()
		posts = append(posts, p)
	}
	log.Debug().Int("num_posts", len(posts)).Msg("GetReplies")
	return posts, true
}

func UpdateReplyCount(db *sql.DB, post int64) bool {
	log.Debug().Msg("UpdateReplyCount")
	stmt, err := db.Prepare("UPDATE post_posts SET num_replies = num_replies + 1 WHERE post_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(post)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}
