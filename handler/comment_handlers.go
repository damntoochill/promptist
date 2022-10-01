package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/art"
	"github.com/promptist/web/comment"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// CommentNewAction updates the profile
func CommentNewAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyForm := c.PostForm("body")
		artUUIDForm := c.PostForm("art-uuid")
		if bodyForm == "" || artUUIDForm == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		piece, ok := art.GetPieceByImageUUID(db, artUUIDForm)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		_, ok = comment.NewComment(db, bodyForm, piece, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/art/"+artUUIDForm)
	}
}

// ProfileCommentNewAction updates the profile
func ProfileCommentNewAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyForm := c.PostForm("body")
		profileUserIDStr := c.PostForm("profile-user-id")
		if bodyForm == "" || profileUserIDStr == "" {
			log.Debug().Msg("body or profile id empty")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		profileUserID, err := strconv.ParseInt(profileUserIDStr, 10, 64)
		if err != nil {
			log.Debug().Msg("unable to convert profile id string")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		profile, ok := profile.GetProfile(db, profileUserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		_, ok = comment.NewProfileComment(db, bodyForm, profile.UserID, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/people/"+profile.Username)
	}
}
