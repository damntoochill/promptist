package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/art"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// LikesPage foo
func LikesPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		usernameParam := c.Param("username")
		p, ok := profile.GetProfileByUsername(db, usernameParam)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		if p.UserID == 0 {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}

		pieces, ok := art.GetPiecesByLike(db, p.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "likes.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
			"Pieces":      pieces,
		})
	}
}

// LikeAction foo
func LikeAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		pieceUIIDStr := c.Query("piece")
		if pieceUIIDStr == "" {
			log.Debug().Msg("piece query is empty")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		piece, ok := art.GetPieceByImageUUID(db, pieceUIIDStr)
		if !ok {
			log.Debug().Msg("unable to get art piece")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		like, ok := art.GetLike(db, piece.ID, us)
		if !ok {
			log.Debug().Msg("unable to get like")

			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if like.ID == 0 {
			// new
			art.NewLike(db, piece.ID, us)
		} else {
			// unlike
			art.Unlike(db, piece.ID, us)
		}

		c.Redirect(http.StatusFound, "/art/"+piece.ImageUUID)
	}
}
