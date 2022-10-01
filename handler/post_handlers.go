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

// PostNewAction updates the profile
func PostNewAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyForm := c.PostForm("body")
		forumIDStr := c.PostForm("forum")
		var f post.Forum
		if bodyForm == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if forumIDStr != "" {
			forumID, err := strconv.ParseInt(forumIDStr, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg("Unable to parse forum ID to int")
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
			var ok bool
			f, ok = post.GetForum(db, forumID)
			if !ok {
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
		}

		_, ok := post.NewPost(db, bodyForm, f.ID, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if f.ID == 0 {
			c.Redirect(http.StatusFound, "/")
		} else {
			c.Redirect(http.StatusFound, "/forums/"+forumIDStr)
		}

	}
}

// PostPage for the user to reply to a standard post
func PostPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		postIDStr := c.Query("post")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msg("unable to parse post ID")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if postID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		p, ok := post.GetPost(db, postID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		profile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		replies, ok := post.GetReplies(db, postID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		log.Debug().Msg("replies: " + strconv.Itoa(len(replies)))

		c.HTML(http.StatusOK, "post.tmpl", gin.H{
			"UserSession": us,
			"Profile":     profile,
			"Post":        p,
			"Replies":     replies,
		})
	}
}

// ReplyNewPage for the user to reply to a standard post
func ReplyNewPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		postIDStr := c.Query("post")
		postID, _ := strconv.ParseInt(postIDStr, 10, 64)

		if postID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		post, ok := post.GetPost(db, postID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		profile, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "reply_new.tmpl", gin.H{
			"UserSession": us,
			"Profile":     profile,
			"Post":        post,
		})
	}
}

// ReplyNewAction for the user to reply to a standard post
func ReplyNewAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Msg("ReplyNewAction")
		bodyForm := c.PostForm("body")
		if bodyForm == "" {
			log.Debug().Msg("ReplyNewAction: bodyForm is empty")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		postIDStr := c.PostForm("post")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			log.Debug().Err(err).Msg("unable to parse post ID")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if postID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		_, ok := post.NewReply(db, bodyForm, postID, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/post/view?post="+postIDStr)
	}
}
