package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/follow"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
)

// FollowAction will follow somebody
func FollowAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		usernameQuery := c.Query("u")
		r := c.Query("r")
		id := c.Query("id")
		leader, ok := profile.GetProfileByUsername(db, usernameQuery)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if leader.UserID == us.UserID {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		// Check if follow relationship exists

		followExists := follow.RelationshipExists(db, leader.UserID, us.UserID)

		if !followExists {
			_, ok := follow.NewRelationship(db, leader.UserID, us.UserID)
			if !ok {
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
		}
		redirect := "/people/" + usernameQuery
		switch r {
		case "piece":
			redirect = "/art/" + id
		case "profile":
			redirect = "/people/" + id
		}
		c.Redirect(http.StatusFound, redirect)
	}
}

// UnfollowAction will unfollow somebody
func UnfollowAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		usernameQuery := c.Query("u")
		r := c.Query("r")
		id := c.Query("id")
		leader, ok := profile.GetProfileByUsername(db, usernameQuery)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		relationship, ok := follow.GetRelationship(db, leader.UserID, us.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if relationship.ID == 0 {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		follow.Unfollow(db, relationship.ID)

		redirect := "/people/" + usernameQuery
		switch r {
		case "piece":
			redirect = "/art/" + id
		case "profile":
			redirect = "/people/" + id
		}
		c.Redirect(http.StatusFound, redirect)
	}
}

// FollowingPage is the omega
func FollowingPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		usernameForm := c.Param("username")
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

		profiles, ok := follow.GetFollowing(db, p.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "follow.tmpl", gin.H{
			"Profile":     p,
			"UserSession": us,
			"Profiles":    profiles,
		})
	}
}

// FollowersPage is the omega
func FollowersPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		usernameForm := c.Param("username")
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

		profiles, ok := follow.GetFollowers(db, p.UserID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "follow.tmpl", gin.H{
			"Profile":     p,
			"UserSession": us,
			"Profiles":    profiles,
		})
	}
}
