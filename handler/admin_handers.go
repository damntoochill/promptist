package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/usersession"
)

// AdminDashboard shows all the good info for the admins
func AdminDashboard(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_dashboard.tmpl", gin.H{
			"UserSession": us,
		})
	}
}
