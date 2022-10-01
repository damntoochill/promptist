package post

import (
	"database/sql"

	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// NewPost creates a new post
func NewPost(db *sql.DB, body string, forumID int64, us *usersession.UserSession) (int64, bool) {
	brief, _ := profile.GetBrief(db, us.UserID)

	stmt, err := db.Prepare("INSERT INTO post_posts (body, user_id, username, full_name, avatar, forum_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(body, us.UserID, brief.Username, brief.FullName, brief.PhotoUUID, forumID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	postID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	if forumID != 0 {
		updateForumPostCount(db, forumID)
	}

	return postID, true
}

// GetPost returns the post from the database
func GetPost(db *sql.DB, postID int64) (Post, bool) {
	stmt, err := db.Prepare("SELECT post_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, post_type, parent_ID FROM post_posts WHERE post_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Post{}, false
	}
	rows, err := stmt.Query(postID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Post{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Post{}, false
	}
	var p Post
	err = rows.Scan(&p.ID, &p.Body, &p.CreatedAt, &p.UserID, &p.Username, &p.FullName, &p.Avatar, &p.NumReplies, &p.NumLikes, &p.PostType, &p.ParentID)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Post{}, false
	}
	p.CreatedAtPretty = timeago.New(p.CreatedAt).Format()
	return p, true
}

// GetPosts from the database
func GetPosts(db *sql.DB) ([]Post, bool) {
	q := `
		SELECT post_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, post_type, parent_id
		FROM post_posts
		WHERE post_type = 0
		ORDER BY created_at DESC
		`
	res, err := db.Query(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Post{}, false
	}
	defer res.Close()

	var posts []Post

	for res.Next() {
		var p Post
		err := res.Scan(&p.ID, &p.Body, &p.CreatedAt, &p.UserID, &p.Username, &p.FullName, &p.Avatar, &p.NumReplies, &p.NumLikes, &p.PostType, &p.ParentID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Post{}, false
		}
		p.CreatedAtPretty = timeago.New(p.CreatedAt).Format()
		posts = append(posts, p)
	}
	return posts, true
}

// GetPostsByForum from the database
func GetPostsByForum(db *sql.DB, forumID int64) ([]Post, bool) {
	log.Debug().Msg("GetPostsByForum")
	q := `
		SELECT post_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, post_type, parent_id
		FROM post_posts
		WHERE forum_id = ? 
		AND post_type = 0
		ORDER BY created_at DESC
		`
	res, err := db.Query(q, forumID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Post{}, false
	}
	defer res.Close()

	var posts []Post

	for res.Next() {
		var p Post
		err := res.Scan(&p.ID, &p.Body, &p.CreatedAt, &p.UserID, &p.Username, &p.FullName, &p.Avatar, &p.NumReplies, &p.NumLikes, &p.PostType, &p.ParentID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Post{}, false
		}
		p.CreatedAtPretty = timeago.New(p.CreatedAt).Format()
		posts = append(posts, p)
	}
	return posts, true
}

// GetPostsByThread from the database
func GetPostsByThread(db *sql.DB, threadID int64) ([]Post, bool) {
	log.Debug().Msg("GetPostsByThread")
	q := `
		SELECT post_id, body, created_at, user_id, username, full_name, avatar, num_replies, num_likes, post_type, parent_id
		FROM post_posts
		WHERE thread_id = ? 
		AND post_type = 0
		ORDER BY created_at 
		`
	res, err := db.Query(q, threadID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Post{}, false
	}
	defer res.Close()

	var posts []Post

	for res.Next() {
		var p Post
		err := res.Scan(&p.ID, &p.Body, &p.CreatedAt, &p.UserID, &p.Username, &p.FullName, &p.Avatar, &p.NumReplies, &p.NumLikes, &p.PostType, &p.ParentID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Post{}, false
		}
		p.CreatedAtPretty = timeago.New(p.CreatedAt).Format()
		posts = append(posts, p)
	}
	return posts, true
}

// GetPostTypeString returns the type in string format
// This is more for debugging
func GetPostTypeString(postType int32) string {
	switch postType {
	case 0:
		return "standard post"
	case 1:
		return "standard reply"
	case 2:
		return "art comment"
	}
	return "error"
}
