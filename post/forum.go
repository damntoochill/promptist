package post

import (
	"database/sql"

	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// GetForums from the database
func GetForums(db *sql.DB) ([]Forum, bool) {
	q := `
		SELECT forum_id, name, slug, about, num_posts, num_views, updated_at
		FROM post_forums
		ORDER BY display_order
		`
	res, err := db.Query(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Forum{}, false
	}
	defer res.Close()

	var forums []Forum

	for res.Next() {
		var f Forum
		err := res.Scan(&f.ID, &f.Name, &f.Slug, &f.About, &f.NumPosts, &f.NumViews, &f.UpdatedAt)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Forum{}, false
		}
		f.UpdatedAtPretty = timeago.New(f.UpdatedAt).Format()
		forums = append(forums, f)
	}
	return forums, true
}

// GetForum returns foo
func GetForum(db *sql.DB, forumID int64) (Forum, bool) {
	stmt, err := db.Prepare("SELECT forum_id, name, slug, about, num_posts, num_views, updated_at FROM post_forums WHERE forum_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Forum{}, false
	}
	rows, err := stmt.Query(forumID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Forum{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Forum{}, false
	}
	var f Forum
	err = rows.Scan(&f.ID, &f.Name, &f.Slug, &f.About, &f.NumPosts, &f.NumViews, &f.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Forum{}, false
	}
	f.UpdatedAtPretty = timeago.New(f.UpdatedAt).Format()
	return f, true
}

// GetForum returns foo
func GetForumBySlug(db *sql.DB, slug string) (Forum, bool) {
	stmt, err := db.Prepare("SELECT forum_id, name, slug, about, num_posts, num_views, updated_at FROM post_forums WHERE slug = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Forum{}, false
	}
	rows, err := stmt.Query(slug)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Forum{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Forum{}, false
	}
	var f Forum
	err = rows.Scan(&f.ID, &f.Name, &f.Slug, &f.About, &f.NumPosts, &f.NumViews, &f.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Forum{}, false
	}
	f.UpdatedAtPretty = timeago.New(f.UpdatedAt).Format()
	return f, true
}

// GetThreads from the database
func GetThreads(db *sql.DB, forumID int64) ([]Thread, bool) {
	q := `
		SELECT id, name, user_id, username, full_name, avatar, last_user_id, last_username, last_full_name, last_avatar, num_replies, is_sticky, is_locked, forum_id, created_at, updated_at, last_message
		FROM post_threads
		WHERE forum_id = ?
		ORDER BY updated_at DESC
		`
	res, err := db.Query(q, forumID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Thread{}, false
	}
	defer res.Close()

	var threads []Thread

	for res.Next() {
		var t Thread
		err := res.Scan(&t.ID, &t.Name, &t.UserID, &t.Username, &t.FullName, &t.Avatar, &t.LastUserID, &t.LastUsername, &t.LastFullName, &t.LastAvatar, &t.NumReplies, &t.IsSticky, &t.IsLocked, &t.ForumID, &t.CreatedAt, &t.UpdatedAt, &t.LastMessage)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Thread{}, false
		}
		t.UpdatedAtPretty = timeago.New(t.UpdatedAt).Format()
		t.CreatedAtPretty = timeago.New(t.CreatedAt).Format()
		threads = append(threads, t)
	}
	return threads, true
}

// GetThread returns foo
func GetThread(db *sql.DB, threadID int64) (Thread, bool) {
	stmt, err := db.Prepare("SELECT id, name, user_id, username, full_name, avatar, last_user_id, last_username, last_full_name, last_avatar, num_replies, is_sticky, is_locked, forum_id, created_at, updated_at, last_message FROM post_threads WHERE id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Thread{}, false
	}
	rows, err := stmt.Query(threadID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Thread{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Thread{}, false
	}
	var t Thread
	err = rows.Scan(&t.ID, &t.Name, &t.UserID, &t.Username, &t.FullName, &t.Avatar, &t.LastUserID, &t.LastUsername, &t.LastFullName, &t.LastAvatar, &t.NumReplies, &t.IsSticky, &t.IsLocked, &t.ForumID, &t.CreatedAt, &t.UpdatedAt, &t.LastMessage)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Thread{}, false
	}
	t.UpdatedAtPretty = timeago.New(t.UpdatedAt).Format()
	t.CreatedAtPretty = timeago.New(t.CreatedAt).Format()
	return t, true
}

// NewThread creates a new post
func NewThread(db *sql.DB, body string, title string, forumID int64, us *usersession.UserSession) (int64, bool) {
	forum, ok := GetForum(db, forumID)
	if !ok {
		return 0, false
	}

	brief, ok := profile.GetBrief(db, us.UserID)
	if !ok {
		log.Debug().Msg("unable to get brief")
		return 0, false
	}

	stmt, err := db.Prepare("INSERT INTO post_threads (name, user_id, username, full_name, avatar, last_user_id, last_username, last_full_name, last_avatar, forum_id, last_message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(title, us.UserID, us.Username, brief.FullName, brief.PhotoUUID, us.UserID, us.Username, brief.FullName, brief.PhotoUUID, forumID, body)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	threadID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	stmt, err = db.Prepare("INSERT INTO post_posts (body, user_id, username, full_name, avatar, forum_id, thread_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err = stmt.Exec(body, us.UserID, brief.Username, brief.FullName, brief.PhotoUUID, forum.ID, threadID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	_, err = res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	//updateForumPostCount(db, thread.ForumID)

	return threadID, true
}

// PostReply creates a new post
func PostReply(db *sql.DB, body string, threadID int64, us *usersession.UserSession) (int64, bool) {
	thread, ok := GetThread(db, threadID)
	if !ok {
		return 0, false
	}

	brief, ok := profile.GetBrief(db, us.UserID)
	if !ok {
		log.Debug().Msg("unable to get brief")
		return 0, false
	}

	stmt, err := db.Prepare("INSERT INTO post_posts (body, user_id, username, full_name, avatar, forum_id, thread_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(body, us.UserID, brief.Username, brief.FullName, brief.PhotoUUID, thread.ForumID, threadID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	postID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	updateForumPostCount(db, thread.ForumID)

	return postID, true
}

func updateForumPostCount(db *sql.DB, forumID int64) bool {
	log.Debug().Msg("UpdateReplyCount")
	stmt, err := db.Prepare("UPDATE post_forums SET num_posts = num_posts + 1 WHERE forum_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(forumID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}

func updateForumViewCount(db *sql.DB, forumID int64) bool {
	log.Debug().Msg("UpdateReplyCount")
	stmt, err := db.Prepare("UPDATE post_forums SET num_views = num_views + 1 WHERE forum_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(forumID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}
