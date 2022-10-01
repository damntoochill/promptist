package message

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/auth"
	"github.com/promptist/web/email"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// NewMessage creates a new message
func NewMessage(db *sql.DB, body string, chatID int64, fromUserID int64, fromName string) (int64, bool) {

	stmt, err := db.Prepare("INSERT INTO message_messages (user_id, body, chat_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(fromUserID, body, chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	messageID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	lastMessage := body
	if len(body) > 80 {
		lastMessage = body[0:75]
	}

	query := `UPDATE message_chats SET last_message=?, updated_at=NOW() WHERE id=?`
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	_, err = stmt.Exec(lastMessage, chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}

	// update is_read

	stmt, err = db.Prepare("UPDATE message_user_chats SET is_read=0 WHERE chat_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	_, err = stmt.Exec(chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return 0, false
	}

	// update unread_chats for other users
	chatUsers, ok := GetChatUsers(db, chatID)
	if !ok {
		log.Error().Msg("unable to get chat users")
	}

	delete(chatUsers, fromUserID)

	for _, mem := range chatUsers {
		query := `UPDATE profile_profiles SET unread_chats=1 WHERE user_id=?`
		stmt, err = db.Prepare(query)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare query")
			return 0, false
		}
		_, err = stmt.Exec(mem.UserID)
		if err != nil {
			log.Error().Err(err).Msg("unable to execute statement")
			return 0, false
		}

		// send email notif
		userEmailAddress, err := auth.Email(db, mem.UserID)
		if err == nil {
			subject := "New message from " + fromName
			chatIDStr := strconv.FormatInt(chatID, 10)
			emailBody := mem.Username + "! You've got a new message on Promptist! You can view it here: " + os.Getenv("WEBSITE_HOST") + "/chat?chat=" + chatIDStr

			err = email.Email(subject, emailBody, userEmailAddress)
			if err != nil {
				log.Error().Err(err).Msg("unable to send email")
			}
		}
	}

	return messageID, true
}

// CreateChat creates a new chat
func CreateChat(db *sql.DB, lastMessage string, userIDs ...int64) (int64, bool) {

	lookupID := MakeLookupID(userIDs)

	briefs := make(map[int64]profile.Brief)
	for _, userID := range userIDs {
		brief, ok := profile.GetBrief(db, userID)
		if !ok {
			log.Debug().Msg("unable to get brief")
			return 0, false
		}
		briefs[userID] = brief
	}

	bz, err := json.Marshal(briefs)
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO message_chats (last_message, member_lookup, members) VALUES (?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(lastMessage, lookupID, bz)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	chatID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	for _, userID := range userIDs {
		CreateUserChat(db, userID, chatID)
	}

	return chatID, true
}

// CreateUserChat creates a new chat
func CreateUserChat(db *sql.DB, userID int64, chatID int64) (int64, bool) {

	stmt, err := db.Prepare("INSERT INTO message_user_chats (user_id, chat_id) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	res, err := stmt.Exec(userID, chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	userChatID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	return userChatID, true
}

// GetChatID returns the ChatID for a given set of users
func GetChatID(db *sql.DB, lookupID string) (int64, bool) {
	log.Debug().Str("lookupID", lookupID).Msg("GetChatID")

	stmt, err := db.Prepare("SELECT id FROM message_chats WHERE member_lookup = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	rows, err := stmt.Query(lookupID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return 0, false
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, false
	}
	var chatID int64
	err = rows.Scan(&chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return 0, false
	}

	return chatID, true
}

func GetLookupID(userIDs ...int) string {
	sort.Ints(userIDs)

	var usersStr []string
	for _, i := range userIDs {
		usersStr = append(usersStr, strconv.Itoa(i))
	}
	return strings.Join(usersStr, ",")
}

func MakeLookupID(userIDs []int64) string {
	sort.Slice(userIDs, func(i, j int) bool { return userIDs[i] < userIDs[j] })

	var usersStr []string
	for _, i := range userIDs {
		usersStr = append(usersStr, strconv.FormatInt(i, 10))
	}
	return strings.Join(usersStr, ",")
}

// GetChat returns the chat from the database
func GetChat(db *sql.DB, chatID int64, myUserID int64) (Chat, bool) {
	stmt, err := db.Prepare("SELECT id, created_at, updated_at, last_message, members FROM message_chats WHERE id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Chat{}, false
	}
	rows, err := stmt.Query(chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Chat{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Chat{}, false
	}
	var c Chat
	var members []byte
	err = rows.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.LastMessage, &members)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Chat{}, false
	}

	err = json.Unmarshal(members, &c.Members)
	if err != nil {
		fmt.Println(err)
		return Chat{}, false
	}
	delete(c.Members, myUserID)
	for _, mem := range c.Members {
		c.Recipient = mem
	}

	if len(c.Members) > 1 {
		c.IsGroup = true
	} else {
		c.IsGroup = false
	}
	c.CreatedAtPretty = timeago.New(c.CreatedAt).Format()
	c.UpdatedAtPretty = timeago.New(c.UpdatedAt).Format()

	// update is_read

	stmt, err = db.Prepare("UPDATE message_user_chats SET is_read=1 WHERE user_id=? AND chat_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Chat{}, false
	}
	_, err = stmt.Exec(myUserID, chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Chat{}, false
	}
	return c, true
}

// GetMessages from the database
func GetMessages(db *sql.DB, chatID int64) ([]Message, bool) {
	log.Debug().Msg("GetMessages")

	chatUsers, ok := GetChatUsers(db, chatID)
	if !ok {
		return []Message{}, false
	}

	q := `
		SELECT id, user_id, body, created_at, chat_id
		FROM message_messages
		WHERE chat_id = ? 
		ORDER BY created_at
		`
	res, err := db.Query(q, chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Message{}, false
	}
	defer res.Close()

	var messages []Message

	for res.Next() {
		var m Message
		err := res.Scan(&m.ID, &m.UserID, &m.Body, &m.CreatedAt, &m.ChatID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Message{}, false
		}
		m.CreatedAtPretty = timeago.New(m.CreatedAt).Format()
		m.Username = chatUsers[m.UserID].Username
		m.FullName = chatUsers[m.UserID].FullName
		m.Avatar = chatUsers[m.UserID].PhotoUUID

		messages = append(messages, m)
	}
	return messages, true
}

// GetChats from the database
func GetChats(db *sql.DB, userID int64) ([]Chat, bool) {
	log.Debug().Msg("GetChats")

	q := `
		SELECT mc.id, mc.created_at, mc.updated_at, last_message, is_read, members
		FROM message_user_chats muc
		JOIN message_chats mc ON mc.id = muc.chat_id
		WHERE muc.user_id = ? 
		ORDER BY mc.updated_at DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Chat{}, false
	}
	defer res.Close()

	var chats []Chat

	for res.Next() {
		var c Chat
		var members []byte
		err := res.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.LastMessage, &c.IsRead, &members)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Chat{}, false
		}
		log.Debug().Msgf("members: %s", members)
		err = json.Unmarshal(members, &c.Members)
		if err != nil {
			fmt.Println(err)
			return nil, false
		}

		delete(c.Members, userID)

		for _, mem := range c.Members {
			c.Recipient = mem
		}

		if len(c.Members) > 1 {
			c.IsGroup = true
		} else {
			c.IsGroup = false
		}
		c.CreatedAtPretty = timeago.New(c.CreatedAt).Format()
		c.UpdatedAtPretty = timeago.New(c.UpdatedAt).Format()
		chats = append(chats, c)
	}
	return chats, true
}

// GetChatUsers from the database
func GetChatUsers(db *sql.DB, chatID int64) (map[int64]profile.Brief, bool) {
	log.Debug().Msg("GetChatUsers")
	q := `
		SELECT pp.id, muc.user_id, username, full_name, photo_uuid
		FROM message_user_chats muc
		JOIN profile_profiles pp ON pp.user_id = muc.user_id
		WHERE chat_id = ? 
		`
	res, err := db.Query(q, chatID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return nil, false
	}
	defer res.Close()

	var briefs map[int64]profile.Brief
	briefs = make(map[int64]profile.Brief)
	for res.Next() {

		var b profile.Brief
		err := res.Scan(&b.ID, &b.UserID, &b.Username, &b.FullName, &b.PhotoUUID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return nil, false
		}
		briefs[b.UserID] = b
	}

	return briefs, true
}

func ValidateChatUser(db *sql.DB, chatID int64, userID int64) bool {
	log.Debug().Msg("ValidateChatUser")
	q := `
		SELECT id
		FROM message_user_chats
		WHERE chat_id = ? AND user_id = ?
		`
	res, err := db.Query(q, chatID, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return false
	}
	defer res.Close()

	if !res.Next() {
		return false
	}
	return true
}

// UnreadChats returns the post from the database
func UnreadChats(db *sql.DB, userID int64) (bool, bool) {
	stmt, err := db.Prepare("SELECT unread_chats FROM profile_profiles WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false, false
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false, false
	}
	defer rows.Close()
	if !rows.Next() {
		return false, false
	}
	var unread bool
	err = rows.Scan(&unread)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return false, false
	}
	return unread, true
}

func CheckForUnread(c *gin.Context, db *sql.DB, us *usersession.UserSession) bool {
	log.Debug().Msg("CheckForUnread")

	unread, ok := UnreadChats(db, us.UserID)
	if !ok {
		return false
	}

	usersession.Save(c, us.IsAdmin, us.IsAuthenticated, us.IsVerified, us.UserID, us.Username, us.UnreadNotifications, unread)

	return true
}

func SetChatsToRead(db *sql.DB, userID int64) bool {
	log.Debug().Msg("SetChatsToRead")
	stmt, err := db.Prepare("UPDATE profile_profiles SET unread_chats=0 WHERE user_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}
