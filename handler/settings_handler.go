package handler

import (
	"database/sql"
	"net/http"

	"github.com/promptist/web/art"
	"github.com/promptist/web/auth"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// SettingsPage shows the user their settings
func SettingsPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		email, err := auth.Email(db, us.UserID)
		if err != nil {
			log.Error().Err(err).Msg("couldn't get email")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
			"Email":       email,
		})
	}
}

// MyArtPage for managing your art
func MyArtPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := art.GetPiecesByUserID(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.HTML(http.StatusOK, "my_art.tmpl", gin.H{
			"UserSession": us,
			"Pieces":      p,
		})
	}
}
