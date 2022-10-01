package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/promptist/web/email"

	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// InitializeReset the password for the given user ID.
func InitializeReset(db *sql.DB, userID int64) error {
	token := xid.New()
	expiry := time.Now().Add(time.Minute * time.Duration(30))
	stmt, err := db.Prepare("INSERT INTO auth_resets (user_id, token, expiry) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, token, expiry)
	if err != nil {
		return err
	}
	e, err := Email(db, userID)
	if err != nil {
		return err
	}
	err = email.PasswordReset(e, token.String())
	if err != nil {
		return err
	}
	return nil
}

// CheckResetToken checks to see if the reset token is valid, if so
// return user id, otherwise return 0
func CheckResetToken(db *sql.DB, email string, token string) error {
	q := `	SELECT email
			FROM auth_resets
			WHERE token=?`
	res, err := db.Query(q, token)
	if err != nil {
		return err
	}
	defer res.Close()
	var emailInDB string
	for res.Next() {
		err := res.Scan(&emailInDB)
		if err != nil {
			return err
		}
	}
	if email != emailInDB {
		return errors.New("bad token")
	}
	return nil
}

// Reset checks to see if the reset token is valid, if so
// return user id, otherwise return 0
func Reset(db *sql.DB, email string, token string, password string) error {
	err := CheckResetToken(db, email, token)
	if err != nil {
		log.Error().Err(err).Msg("bad token")
		return err
	}
	if len(password) < 8 || len(password) > 254 {
		return errors.New("password needs to be between 8 and 254 characters")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error().Err(err).Msg("unable to generate password hash")
		return err
	}
	hashedPassword := string(bytes)
	q := `	UPDATE auth_users 
			SET password_hash=?
			WHERE email=?`
	stmt, err := db.Prepare(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare update statement")
		return err
	}
	_, err = stmt.Exec(hashedPassword, email)
	if err != nil {
		log.Error().Err(err).Msg("unable to update password")
		return err
	}
	return nil
}
