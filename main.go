// Package main starts up the machine
package main

import (
	"database/sql"
	"os"

	"github.com/promptist/web/handler"
	"github.com/promptist/web/middleware"
	"github.com/promptist/web/usersession"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// Logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting ZeroLog logger")

	// Config Management
	err := checkConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("missing config variable")
	}

	// DB
	connString := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_DATABASE") + "?parseTime=true"

	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open database connection")
	}
	defer db.Close()

	// Gin init
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "./public")

	// Sessions
	sessionKey := os.Getenv("SESSION_KEY")
	store := cookie.NewStore([]byte(sessionKey))
	r.Use(sessions.Sessions("mysession", store))
	us := &usersession.UserSession{}

	// Middleware
	r.Use(middleware.PrintSession(us))
	r.Use(middleware.LoadUserSession(us))
	r.Use(middleware.PrintSession(us))

	/////////
	// Routes

	// Home
	r.GET("/", handler.HomePage(db, us))

	// Forums
	r.GET("/forums", handler.ForumsPage(db, us))
	r.GET("/forums/:slug", handler.ForumPage(db, us))
	r.GET("/thread", handler.ThreadPage(db, us))

	// Royales
	r.GET("/royale", handler.RoyalesPage(db, us))

	// Info
	r.GET("/faq", handler.FaqPage(us))
	r.GET("/what", handler.WhatIsThisPage(us))
	r.GET("/fresh", handler.FreshPage(us))
	r.GET("/rules", handler.RulesPage(us))
	r.GET("/pioneers", handler.PioneerPage(db, us))
	r.GET("/pro", handler.ProPage(us))

	// People
	r.GET("/people/:username", handler.ProfilePage(db, us))
	r.GET("/people/", handler.ProfileListPage(db, us))

	// Art
	r.GET("/tag/:name", handler.TagPage(db, us))
	r.GET("/art/:imageUUID", handler.PiecePage(db, us))
	r.GET("/search", handler.SearchPage(db, us))
	r.GET("/program/:slug", handler.ProgramPage(db, us))

	// Post
	r.GET("/post/view", handler.PostPage(db, us))

	// Auth
	r.GET("/join", handler.JoinPage())
	r.POST("/join", handler.JoinAction(db))
	r.GET("/login", handler.LoginPage())
	r.POST("/login", handler.LoginAction(db))
	r.GET("/logout", handler.LogoutAction())
	r.GET("/forgot-password", handler.ForgotPasswordPage())
	r.POST("/forgot-password", handler.ForgotPasswordAction(db))
	r.GET("/reset-password", handler.PasswordResetPage(db))
	r.POST("/reset-password", handler.PasswordResetAction(db))
	r.GET("/verify", handler.VerifyAction(db, us))
	r.GET("/resend-verification", handler.ResendVerificationAction(db, us))

	auth := r.Group("/")
	auth.Use(middleware.Protect(us))
	{
		// Art
		auth.GET("/upload", handler.UploadPage(db, us))
		auth.POST("/upload", handler.UploadAction(db, us))
		auth.GET("/edit/:imageUUID", handler.EditPage(db, us))
		auth.POST("/edit", handler.EditAction(db, us))
		auth.GET("/delete", handler.DeleteArtAction(db, us))

		// Art Like
		auth.GET("/like/new", handler.LikeAction(db, us))

		// Follow
		auth.GET("/follow", handler.FollowAction(db, us))
		auth.GET("/unfollow", handler.UnfollowAction(db, us))

		// ACCOUNT
		auth.GET("/account", handler.SettingsPage(db, us))
		auth.GET("/account/my-art", handler.MyArtPage(db, us))
		auth.GET("/account/profile-photo", handler.ProfilePhotoEditPage(db, us))
		auth.POST("/account/profile-photo", handler.ProfilePhotoEditAction(db, us))

		// Post
		auth.POST("/post/new", handler.PostNewAction(db, us))
		auth.GET("/reply/new", handler.ReplyNewPage(db, us))
		auth.POST("/reply/new", handler.ReplyNewAction(db, us))

		// Forums
		auth.GET("/thread/new", handler.ThreadNewPage(db, us))
		auth.POST("/thread/new", handler.ThreadNewAction(db, us))
		auth.GET("/thread/reply", handler.ThreadReplyPage(db, us))
		auth.POST("/thread/reply", handler.ThreadReplyAction(db, us))

		// Comment
		auth.POST("/comment/new", handler.CommentNewAction(db, us))
		auth.POST("/profilecomment/new", handler.ProfileCommentNewAction(db, us))

		// Profile
		auth.GET("/account/profile", handler.ProfileEditPage(db, us))
		auth.POST("/account/profile", handler.ProfileEditAction(db, us))

		// Collection management
		auth.GET("/collection/new", handler.CollectionNewPage(db, us))
		auth.POST("/collection/new", handler.CollectionNewAction(db, us))
		auth.GET("/collection/prompt", handler.CollectionsPromptPage(db, us))
		auth.GET("/collection/save", handler.SaveToCollectionAction(db, us))
		auth.GET("/collection/edit", handler.CollectionEditPage(db, us))
		auth.POST("/collection/edit", handler.CollectionEditAction(db, us))
		auth.GET("/collection/delete", handler.CollectionDeleteAction(db, us))
		auth.GET("/collection/delete-piece", handler.CollectionPieceDeleteAction(db, us))

		// Notifications
		auth.GET("/notifications", handler.NotificationsPage(db, us))

		// Chat
		auth.GET("/chat/list", handler.ChatListPage(db, us))
		auth.GET("/chat/direct", handler.ProfileMessageButtonAction(db, us))
		auth.POST("/chat/send", handler.SendMessageAction(db, us))
		auth.GET("/chat", handler.ChatPage(db, us))

		// New Player
		auth.GET("/a-new-player", handler.NewPlayerPage(db, us))

		// Royales
		r.GET("/royale/select", handler.RoyaleSelectPage(db, us))
		r.GET("/royale/enter", handler.RoyaleEnterPage(db, us))
		r.GET("/royale/judge", handler.RoyaleAlertAction(db, us))
	}

	admin := r.Group("/admin")
	admin.Use(middleware.Admin(us))
	{
		admin.GET("/", handler.AdminDashboard(db, us))
	}

	// Royales
	r.GET("/royale/:slug", handler.RoyalePage(db, us))

	// Likes
	r.GET("/:username/likes", handler.LikesPage(db, us))

	// Collections
	r.GET("/:username/collections", handler.CollectionsPage(db, us))
	r.GET("/:username/collections/:id", handler.CollectionPage(db, us))

	// Following
	r.GET("/:username/followers", handler.FollowersPage(db, us))
	r.GET("/:username/following", handler.FollowingPage(db, us))

	r.GET("/:username", handler.ProfilePage(db, us))

	r.Run(":" + os.Getenv("PORT"))
}
