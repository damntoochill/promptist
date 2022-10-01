package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/art"
	"github.com/promptist/web/comment"
	"github.com/promptist/web/follow"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/promptist/web/utils"
	"github.com/rs/zerolog/log"
)

// ProfileEditPage allows the user to edit their profile
func ProfileEditPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please tell customer support"},
			})
			return
		}
		c.HTML(http.StatusOK, "profile_edit.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
		})
	}
}

// ProfileEditAction updates the profile
func ProfileEditAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameForm := c.PostForm("username")
		fullNameForm := c.PostForm("full-name")
		bioForm := c.PostForm("bio")
		locationForm := c.PostForm("location")

		bio := utils.NewNullString(bioForm)
		location := utils.NewNullString(locationForm)

		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}

		problems, ok := profile.ValidateProfile(usernameForm, fullNameForm)
		if !ok {
			c.HTML(http.StatusOK, "profile_edit.tmpl", gin.H{
				"UserSession": us,
				"Profile":     p,
				"Problems":    problems,
			})
		}

		p.Username = usernameForm
		p.FullName = fullNameForm
		p.Bio = bio
		p.Location = location
		_, err := profile.UpdateProfile(db, p)
		if err != nil {
			log.Error().Err(err).Msg("unable to save profile")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		c.Redirect(http.StatusFound, "/people/"+usernameForm)
	}
}

// ProfileListPage shows the profile page for a user
func ProfileListPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := profile.GetProfiles(db)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		c.HTML(http.StatusOK, "profile_list.tmpl", gin.H{
			"Profiles":    p,
			"UserSession": us,
		})
	}
}

// ProfilePage shows the profile page for a user
func ProfilePage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameForm := c.Param("username")
		var myBrief profile.Brief

		p, ok := profile.GetProfileByUsername(db, usernameForm)
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
		pieces, _ := art.GetPiecesByUserID(db, p.UserID)

		var fo int32
		if us.IsAuthenticated {
			fo = follow.GetFollowOption(db, p.UserID, us.UserID)
		}

		if us.IsAuthenticated {
			myBrief, _ = profile.GetBrief(db, us.UserID)
		}

		comments, _ := comment.GetProfileComments(db, p.UserID)

		c.HTML(http.StatusOK, "profile.tmpl", gin.H{
			"Profile":      p,
			"UserSession":  us,
			"Pieces":       pieces,
			"FollowOption": fo,
			"MyBrief":      myBrief,
			"Comments":     comments,
			"Title":        p.FullName + " | Promptist",
		})
	}
}
