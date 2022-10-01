package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/message"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// ChatListPage for the user to reply to a standard message
func ChatListPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		chats, ok := message.GetChats(db, us.UserID)
		if !ok {
			log.Debug().Msg("unable to get chats")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		myBrief, ok := profile.GetBrief(db, us.UserID)
		if !ok {
			log.Debug().Msg("unable to get my brief")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		message.SetChatsToRead(db, us.UserID)
		usersession.Save(c, us.IsAdmin, us.IsAuthenticated, us.IsVerified, us.UserID, us.Username, us.UnreadNotifications, false)

		c.HTML(http.StatusOK, "chat_list.tmpl", gin.H{
			"UserSession": us,
			"Chats":       chats,
			"MyBrief":     myBrief,
		})
	}
}

// ChatPage for the user to reply to a standard message
func ChatPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatIDStr := c.Query("chat")
		if chatIDStr == "" {
			log.Debug().Msg("no chat id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Debug().Msg("unable to parse chat id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if chatID == 0 {
			log.Debug().Msg("chat id is 0")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		ok := message.ValidateChatUser(db, chatID, us.UserID)
		if !ok {
			log.Debug().Msg("user is not a member of chat")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		chat, ok := message.GetChat(db, chatID, us.UserID)
		if !ok {
			log.Debug().Msg("unable to get chat")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		messages, ok := message.GetMessages(db, chatID)
		if !ok {
			log.Debug().Msg("unable to get messages")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		myBrief, ok := profile.GetBrief(db, us.UserID)
		if !ok {
			log.Debug().Msg("unable to get my brief")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "chat.tmpl", gin.H{
			"UserSession": us,
			"Chat":        chat,
			"Messages":    messages,
			"MyBrief":     myBrief,
		})
	}
}

// ProfileMessageButtonAction for when a user clicks the message button on a profile
func ProfileMessageButtonAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user")
		if userIDStr == "" {
			log.Debug().Msg("no user id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Debug().Msg("unable to parse user id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if userID == 0 {
			log.Debug().Msg("user id is 0")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		userIDs := []int64{us.UserID, userID}
		lookupID := message.MakeLookupID(userIDs)
		var chatID int64
		chatID, _ = message.GetChatID(db, lookupID)
		if chatID == 0 {
			newChatID, ok := message.CreateChat(db, "", us.UserID, userID)
			if !ok {
				log.Debug().Msg("unable to create chat")
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
			chatID = newChatID
		}
		c.Redirect(http.StatusFound, "/chat?chat="+strconv.FormatInt(chatID, 10))
	}
}

// SendMessageAction updates the profile
func SendMessageAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyForm := c.PostForm("body")
		chatIDStr := c.PostForm("chat")
		if bodyForm == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if chatIDStr == "" {
			log.Debug().Msg("no chat id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Debug().Msg("unable to parse chat id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if chatID == 0 {
			log.Debug().Msg("chat id is 0")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		ok := message.ValidateChatUser(db, chatID, us.UserID)
		if !ok {
			log.Debug().Msg("user is not a member of chat")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		_, ok = message.NewMessage(db, bodyForm, chatID, us.UserID, us.Username)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/chat?chat="+strconv.FormatInt(chatID, 10))
	}
}
