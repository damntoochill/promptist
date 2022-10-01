package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/usersession"
)

// NewPlayerPage provides answers
func NewPlayerPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_player.tmpl", gin.H{
			"UserSession": us,
			"Title":       "A New Player Has Entered the Game",
		})
	}
}
