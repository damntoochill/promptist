package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/notification"
	"github.com/promptist/web/usersession"
)

// NotificationsPage shows the profile page for a user
func NotificationsPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		notifications, ok := notification.GetNotifications(db, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		notification.SetNotificationsToRead(db, us.UserID)
		usersession.Save(c, us.IsAdmin, us.IsAuthenticated, us.IsVerified, us.UserID, us.Username, false, us.UnreadChats)

		c.HTML(http.StatusOK, "notifications.tmpl", gin.H{
			"UserSession":   us,
			"Notifications": notifications,
		})
	}
}
