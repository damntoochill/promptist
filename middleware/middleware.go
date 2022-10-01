package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// Admin will abort anyone that is not an admin
func Admin(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !us.IsAdmin {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			c.Abort()
			return
		}
	}
}

// Protect will abort anything that doesn't have admin or auth
func Protect(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !us.IsAuthenticated {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			c.Abort()
			return
		}
	}
}

// PrintSession is middleware that will print the current session info
func PrintSession(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Int64("UserID", us.UserID).Bool("IsAuthenticated", us.IsAuthenticated).Msg("UserSession data")
	}
}

// LoadUserSession will abort anything that doesn't have admin or auth
func LoadUserSession(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		usersession.Load(c, us)
	}
}
