// Package handler holds all the http handlers
package handler

import (
	"database/sql"
	"net/http"

	"github.com/promptist/web/auth"
	"github.com/promptist/web/message"
	"github.com/promptist/web/notification"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ForgotPasswordAction resets the password if everything looks good
func ForgotPasswordAction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.GetUserByEmail(db, c.PostForm("email"))
		if err != nil {
			log.Error().Err(err).Msg("database error getting user id by email")
			c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
				"Error": "there was a database error",
			})
		}
		if userID == 0 {
			log.Info().Err(err).Msg("email not found during forgot password lookup")
			c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
				"Error": "E-mail not found",
			})
		}
		err = auth.InitializeReset(db, userID)
		if err != nil {
			log.Error().Err(err).Msg("unable to initialize reset")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"there was a problem resetting your password"},
			})
			return
		}
		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
			"Success": true,
			"Error":   nil,
		})
	}
}

// ForgotPasswordPage orchestrates the construction of the forgot password
// page that is displayed to a user when they can'r remember their login
// details.
func ForgotPasswordPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
			"Error": "",
		})
	}
}

// JoinPage orchestrates the construction of the join page before
// handing it off to the UI
func JoinPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "join.tmpl", gin.H{
			"ValidationErrors": nil,
		})
	}
}

// JoinAction handles the action of when a user joins
func JoinAction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		errors := auth.RegisterUser(db, c.PostForm("email"), c.PostForm("full-name"), c.PostForm("password"))
		if len(errors) > 0 {
			c.HTML(http.StatusOK, "join.tmpl", gin.H{
				"ValidationErrors": errors,
			})
			return
		}
		userID, err := auth.Login(db, c.PostForm("email"), c.PostForm("password"))
		if err != nil {
			log.Error().Err(err).Msg("error logging in")
			errors := []string{"There was a login issue"}
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}
		p, ok := profile.GetProfile(db, userID)
		if !ok {
			errors := []string{"There was a login issue"}
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}
		usersession.Save(c, false, true, false, userID, p.Username, false, true)

		chatID, ok := message.CreateChat(db, "", 2, userID)
		if !ok {
			log.Error().Err(err).Msg("error creating chat")
		}
		body := "Hi " + p.FullName + "! My name's Tyler. Welcome to the site! Let me know if you have any questions or if you run into any problems. We're just starting to get things off the ground, but I'm excited to see where we can take this place. Thanks for joining!! -tyler"
		_, ok = message.NewMessage(db, body, chatID, 2, "tyler")

		c.Redirect(http.StatusFound, "/")
	}
}

// LoginAction checks to see if we are okay to enter the app
func LoginAction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.Login(db, c.PostForm("email"), c.PostForm("password"))
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"Error": "invalid login",
			})
			return
		}
		if userID == 0 {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"Error": "invalid login",
			})
			return
		}
		isVerified, err := auth.IsVerified(db, userID)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"Error": "unable to verify",
			})
			return
		}
		isAdmin, err := auth.IsAdmin(db, userID)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"Error": "unable to verify",
			})
			return
		}
		p, ok := profile.GetProfile(db, userID)
		if !ok {
			errors := []string{"There was a login issue"}
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}
		unreadNotificatoins, ok := notification.UnreadNotifications(db, userID)
		if !ok {
			log.Debug().Msg("unable to check unread notifications")
		}
		log.Debug().Bool("unread", unreadNotificatoins).Msg("unread notidications")

		unreadChats, ok := message.UnreadChats(db, userID)
		if !ok {
			log.Debug().Msg("unable to check unread notifications")
		}

		usersession.Save(c, isAdmin, true, isVerified, userID, p.Username, unreadNotificatoins, unreadChats)
		if p.IsNew {
			c.Redirect(http.StatusFound, "/a-new-player")
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

// LoginPage orchestrates the construction of the login page before
// handing it off to the UI
func LoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"Error": "",
		})
	}
}

// LogoutAction handles the action of logging out
func LogoutAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		usersession.Clear(c)
		c.Redirect(http.StatusFound, "/login")
	}
}

// PasswordResetAction takes a new password, token, and email, and everything
// look good, updates the password
func PasswordResetAction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.Reset(db, c.PostForm("email"), c.PostForm("token"), c.PostForm("password"))
		if err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"unable to reset password"},
			})
			return
		}
		c.Redirect(http.StatusFound, "/login")
	}
}

// PasswordResetPage is where the user specifies their new password
// authenticated by a token
func PasswordResetPage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.CheckResetToken(db, c.Query("email"), c.Query("token"))
		if err != nil {
			log.Error().Err(err).Msg("error when checking the reset token")
			errors := []string{"There was an error restting your password"}
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}
		c.HTML(http.StatusOK, "password_reset.tmpl", gin.H{
			"Email": c.Query("email"),
			"Token": c.Query("token"),
		})
	}
}

// VerifyAction checks to see if we have a valid email
func VerifyAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := auth.Verify(db, c.Query("email"), c.Query("token"))
		if err != nil {
			log.Error().Err(err).Msg("Email verification error")
			errors := []string{"unable to complete email verficiation"}
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}
		if !ok {
			log.Info().Msg("Email verification error")
			errors := []string{"unable to complete email verficiation"}
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}
		usersession.Save(c, us.IsAdmin, us.IsAuthenticated, true, us.UserID, us.Username, us.UnreadNotifications, us.UnreadChats)
		c.Redirect(http.StatusFound, "/settings")
	}
}

// ResendVerificationAction checks to see if we have a valid email
func ResendVerificationAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.ResendVerification(db, us.UserID)
		if err != nil {
			log.Error().Err(err).Msg("unable to resend verification email")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"unable to send email"},
			})
			return
		}
		c.HTML(http.StatusOK, "message.tmpl", gin.H{
			"Message": "Verification email sent!",
		})
	}
}
