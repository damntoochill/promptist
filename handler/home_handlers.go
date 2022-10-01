package handler

import (
	"database/sql"
	"net/http"

	"github.com/promptist/web/art"
	"github.com/promptist/web/message"
	"github.com/promptist/web/notification"
	"github.com/promptist/web/post"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"

	"github.com/gin-gonic/gin"
)

// HomePage is the omega
func HomePage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		pieces, _ := art.GetPieces(db)
		brief, _ := profile.GetBrief(db, us.UserID)
		posts, _ := post.GetPosts(db)
		notification.CheckForUnread(c, db, us)
		message.CheckForUnread(c, db, us)
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"UserSession": us,
			"Pieces":      pieces,
			"Brief":       brief,
			"Posts":       posts,
		})
	}
}
