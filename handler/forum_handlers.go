package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/post"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// ForumsPage for the user to reply to a standard post
func ForumsPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		forums, ok := post.GetForums(db)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		myProfile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "forums.tmpl", gin.H{
			"UserSession": us,
			"MyProfile":   myProfile,
			"Forums":      forums,
		})
	}
}

// ForumPage for the user to reply to a standard post
func ForumPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		slugParam := c.Param("slug")

		if slugParam == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		f, ok := post.GetForumBySlug(db, slugParam)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		threads, ok := post.GetThreads(db, f.ID)
		if !ok {
			log.Debug().Msg("unable to get threads")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		myProfile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "forum.tmpl", gin.H{
			"UserSession": us,
			"MyProfile":   myProfile,
			"Forum":       f,
			"Threads":     threads,
		})
	}
}

// ThreadPage for the user to reply to a standard post
func ThreadPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		threadIDStr := c.Query("id")
		threadID, err := strconv.ParseInt(threadIDStr, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msg("unable to parse thread ID")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if threadID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		p, ok := post.GetPostsByThread(db, threadID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		myProfile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		thread, ok := post.GetThread(db, threadID)
		if !ok {
			log.Debug().Msg("unable to get thread")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		forum, ok := post.GetForum(db, thread.ForumID)
		if !ok {
			log.Debug().Msg("unable to get forum")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "thread.tmpl", gin.H{
			"UserSession": us,
			"MyProfile":   myProfile,
			"Posts":       p,
			"Thread":      thread,
			"Forum":       forum,
		})
	}
}

// ThreadNewPage for the user to reply to a standard post
func ThreadNewPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		forumIDStr := c.Query("forum")
		forumID, _ := strconv.ParseInt(forumIDStr, 10, 64)

		if forumID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		forum, ok := post.GetForum(db, forumID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		myProfile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "thread_new.tmpl", gin.H{
			"UserSession": us,
			"MyProfile":   myProfile,
			"Forum":       forum,
		})
	}
}

// ThreadNewAction for the user to reply to a standard post
func ThreadNewAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		forumIDStr := c.PostForm("forum")
		body := c.PostForm("body")
		title := c.PostForm("title")

		forumID, _ := strconv.ParseInt(forumIDStr, 10, 64)

		if forumID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if len(body) == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if len(title) == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		threadID, ok := post.NewThread(db, body, title, forumID, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		threadIDStr := strconv.FormatInt(threadID, 10)
		c.Redirect(http.StatusFound, "/thread?id="+threadIDStr)
	}
}

// ThreadReplyPage for the user to reply to a standard post
func ThreadReplyPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		threadIDStr := c.Query("id")
		threadID, _ := strconv.ParseInt(threadIDStr, 10, 64)

		if threadID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		thread, ok := post.GetThread(db, threadID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		forum, ok := post.GetForum(db, thread.ForumID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		myProfile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "thread_reply.tmpl", gin.H{
			"UserSession": us,
			"MyProfile":   myProfile,
			"Forum":       forum,
			"Thread":      thread,
		})
	}
}

// ThreadReplyAction for the user to reply to a standard post
func ThreadReplyAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		threadIDStr := c.PostForm("thread")
		body := c.PostForm("body")

		threadID, _ := strconv.ParseInt(threadIDStr, 10, 64)

		if threadID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if len(body) == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		thread, ok := post.GetThread(db, threadID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		postID, ok := post.PostReply(db, body, thread.ID, us)

		c.Redirect(http.StatusFound, "/thread?id="+threadIDStr+"&post="+strconv.FormatInt(postID, 10))
	}
}
