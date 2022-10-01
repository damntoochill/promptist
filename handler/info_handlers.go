package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
)

// WhatIsThisPage tells what this site is about
func WhatIsThisPage(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "what.tmpl", gin.H{
			"UserSession": us,
		})
	}
}

// FaqPage provides answers
func FaqPage(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "faq.tmpl", gin.H{
			"UserSession": us,
		})
	}
}

// FreshPage provides answers
func FreshPage(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "fresh.tmpl", gin.H{
			"UserSession": us,
		})
	}
}

// RulesPage provides answers
func RulesPage(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "rules.tmpl", gin.H{
			"UserSession": us,
		})
	}
}

// PioneerPage provides answers
func PioneerPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		pioneers, _ := profile.GetPioneers(db)
		c.HTML(http.StatusOK, "pioneer.tmpl", gin.H{
			"UserSession": us,
			"Pioneers":    pioneers,
			"Title":       "The Pioneers",
		})
	}
}

// ProPage provides answers
func ProPage(us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "pro.tmpl", gin.H{
			"UserSession": us,
			"Title":       "Pro",
		})
	}
}
