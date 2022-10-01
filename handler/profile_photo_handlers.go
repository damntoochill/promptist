package handler

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/image"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

// ProfilePhotoEditPage allows the user to edit their profile
func ProfilePhotoEditPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please tell customer support"},
			})
			return
		}
		c.HTML(http.StatusOK, "profile_photo_edit.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
		})
	}
}

// ProfilePhotoEditAction updates the profile
func ProfilePhotoEditAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := xid.New().String()
		path := "images/" + uuid

		log.Debug().Msg("ProfilePhotoEditAction")

		file, err := c.FormFile("file")
		if err != nil {
			log.Error().Err(err).Msg("unable to get file from form ")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if file == nil {
			log.Debug().Msg("file is nil")
			c.HTML(http.StatusNotAcceptable, "404.tmpl", gin.H{})
			return
		}

		err = c.SaveUploadedFile(file, path)
		if err != nil {
			log.Error().Err(err).Msg("unable to save file to destination ")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		ok := image.Upload(path, path, true, 100, 200, 400, 600)
		if !ok {
			log.Error().Err(err).Msg("unable to upload file")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		err = os.Remove(path)
		if err != nil {
			log.Error().Err(err).Msg("unable to delete file")
		}

		ok = profile.UpdateProfilePhoto(db, us.UserID, uuid)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem saving your info, please tell customer support"},
			})
			return
		}
		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please tell customer support"},
			})
			return
		}
		c.HTML(http.StatusOK, "profile_photo_edit.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
		})
	}
}
