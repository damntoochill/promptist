// Package usersession holds the frequenetly used data for the user
package usersession

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Save the user session
func Save(c *gin.Context, isAdmin bool, isAuthenticated bool, isVerified bool, userID int64, username string, unreadNotifications bool, unreadChats bool) {
	session := sessions.Default(c)
	session.Set("IsAdmin", isAdmin)
	session.Set("IsAuthenticated", isAuthenticated)
	session.Set("IsVerified", isVerified)
	session.Set("UserID", userID)
	session.Set("Username", username)
	session.Set("UnreadNotifications", unreadNotifications)
	session.Set("UnreadChats", unreadChats)
	session.Save()
}

// Clear the user session
func Clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// Load the user session
func Load(c *gin.Context, us *UserSession) {
	session := sessions.Default(c)

	userID := session.Get("UserID")
	if userID == nil {
		us.UserID = 0
	} else {
		us.UserID = userID.(int64)
	}

	username := session.Get("Username")
	if username == nil {
		us.Username = ""
	} else {
		us.Username = username.(string)
	}

	isAuthenticated := session.Get("IsAuthenticated")
	if isAuthenticated == nil {
		us.IsAuthenticated = false
	} else {
		us.IsAuthenticated = isAuthenticated.(bool)
	}

	isVerified := session.Get("IsVerified")
	if isVerified == nil {
		us.IsVerified = false
	} else {
		us.IsVerified = isVerified.(bool)
	}

	isAdmin := session.Get("IsAdmin")
	if isAdmin == nil {
		us.IsAdmin = false
	} else {
		us.IsAdmin = isAdmin.(bool)
	}

	unreadNotifications := session.Get("UnreadNotifications")
	if unreadNotifications == nil {
		us.UnreadNotifications = false
	} else {
		us.UnreadNotifications = unreadNotifications.(bool)
	}

	unreadChats := session.Get("UnreadChats")
	if unreadChats == nil {
		us.UnreadChats = false
	} else {
		us.UnreadChats = unreadChats.(bool)
	}
}
