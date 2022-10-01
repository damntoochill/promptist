package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/art"
	"github.com/promptist/web/royale"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// RoyalesPage for the user to reply to a standard post
func RoyalesPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		royales, ok := royale.GetRoyales(db)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "royale_list.tmpl", gin.H{
			"UserSession": us,
			"Royales":     royales,
			"Title":       "Promptist Royales",
		})
	}
}

// RoyalePage for the user to reply to a standard post
func RoyalePage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		slug := c.Param("slug")
		if slug == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		roy, ok := royale.GetRoyaleBySlug(db, slug)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		subs, ok := royale.GetSubmissions(db, roy.ID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		log.Debug().Int("num posts", len(subs)).Msg("num subs")

		amIJudging := royale.AmIJudging(db, roy.ID, us)

		c.HTML(http.StatusOK, "royale.tmpl", gin.H{
			"UserSession": us,
			"Royale":      roy,
			"Submissions": subs,
			"AmIJudging":  amIJudging,
			"Title":       roy.Name + " Royale - Promptist",
		})
	}
}

// RoyaleSelectPage for the user to reply to a standard post
func RoyaleSelectPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		artID := c.Query("art")
		if artID == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		art, ok := art.GetPieceByImageUUID(db, artID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		royales, ok := royale.GetRoyales(db)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "royale_select.tmpl", gin.H{
			"UserSession": us,
			"Royales":     royales,
			"Art":         art,
		})
	}
}

// RoyaleEnterPage for the user to reply to a standard post
func RoyaleEnterPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		good := true
		artID := c.Query("art")
		if artID == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		art, ok := art.GetPieceByImageUUID(db, artID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		slug := c.Query("royale")
		if slug == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		roy, ok := royale.GetRoyaleBySlug(db, slug)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		ok, problem := royale.Submit(db, roy.ID, art.ImageUUID, us)
		if !ok {
			good = false
		}

		c.HTML(http.StatusOK, "royale_enter.tmpl", gin.H{
			"UserSession": us,
			"Royale":      roy,
			"Art":         art,
			"Good":        good,
			"Problem":     problem,
		})
	}
}

// RoyaleAlertAction foo
func RoyaleAlertAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		slug := c.Query("royale")
		if slug == "" {
			log.Debug().Msg("no royale query slug")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		roy, ok := royale.GetRoyaleBySlug(db, slug)
		if !ok {
			log.Debug().Msg("no royale")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		ok, _ = royale.Alert(db, roy.ID, us)
		if !ok {
			log.Debug().Msg("failed to add alert")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.Redirect(http.StatusFound, "/royale/"+slug)
	}
}
