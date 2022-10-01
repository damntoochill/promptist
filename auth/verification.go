package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/promptist/web/email"

	"github.com/rs/xid"
)

// Verify checks to see if the email is authentic
func Verify(db *sql.DB, email string, token string) (bool, error) {
	q := `	SELECT id, email
			FROM auth_verifications
			WHERE token=?
			ORDER BY created_at
			LIMIT 1`
	res, err := db.Query(q, token)
	if err != nil {
		return false, err
	}
	defer res.Close()

	var userID int
	var emailInDB string

	for res.Next() {
		err := res.Scan(&userID, &emailInDB)
		if err != nil {
			return false, err
		}
	}

	if email != emailInDB {
		return false, errors.New("bad token")
	}

	q = `	UPDATE auth_users 
			SET verified=?
			WHERE id=?`
	stmt, err := db.Prepare(q)
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(1, userID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ResendVerification sends an email with the token
func ResendVerification(db *sql.DB, userID int64) error {
	userEmail, err := Email(db, userID)
	if err != nil {
		return err
	}
	token, err := newVerification(db, userID)
	if err != nil {
		return err
	}
	err = email.Verification(userEmail, token)
	if err != nil {
		return err
	}
	return nil
}

// StartVerification sends an email with the token
func StartVerification(db *sql.DB, userID int64) error {
	userEmail, err := Email(db, userID)
	if err != nil {
		return err
	}
	token, err := newVerification(db, userID)
	if err != nil {
		return err
	}
	err = email.Verification(userEmail, token)
	if err != nil {
		return err
	}
	return nil
}

func newVerification(db *sql.DB, userID int64) (string, error) {
	token := xid.New()
	expiry := time.Now().Add(time.Minute * time.Duration(30))
	stmt, err := db.Prepare("INSERT INTO auth_verifications (user_id, token, expiry) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(userID, token, expiry)
	if err != nil {
		return "", nil
	}
	return token.String(), nil
}
