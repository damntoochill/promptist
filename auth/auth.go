package auth

import (
	"database/sql"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/gosimple/slug"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser adds a user to the system
func RegisterUser(db *sql.DB, email string, fullName string, password string) []string {
	log.Debug().Msg("RegisterUser()")
	var validationErrors []string

	if len(email) < 3 || len(email) > 254 {
		validationErrors = append(validationErrors, "email must be between 3 and 254 characters")
	}
	if len(fullName) < 3 || len(fullName) > 254 {
		validationErrors = append(validationErrors, "Name must be between 3 and 254 characters")
	}
	if len(password) < 8 || len(password) > 254 {
		validationErrors = append(validationErrors, "password must be between 8 and 254 characters")
	}
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		validationErrors = append(validationErrors, "invalid email format")
	}
	q := `SELECT email FROM auth_users WHERE email = ?`
	err := db.QueryRow(q, email).Scan(&email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("Unable to query auth table")
			return []string{"database error"}
		}
	} else {
		validationErrors = append(validationErrors, "email already in use")
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hashedPassword := string(bytes)
	usernameSlug := slug.Make(fullName)

	rand.Seed(time.Now().UnixNano())
	randomNum := utils.Random(1000, 9999)
	t := strconv.Itoa(randomNum)
	username := usernameSlug + "-" + t
	stmt, err := db.Prepare("INSERT INTO auth_users (email, password_hash) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("Unable to insert into auth table")
		return []string{"database error"}
	}
	res, err := stmt.Exec(email, hashedPassword)
	if err != nil {
		log.Error().Err(err).Msg("Unable to execute on auth table")
		return []string{"database error"}
	}
	userID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("Could not get user ID from join handler")
	}

	err = StartVerification(db, userID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to send verification email")
		return []string{"unable to send verification email"}
	}

	p, err := profile.NewProfile(userID, username, fullName, sql.NullString{}, sql.NullString{})
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new profile")
		return []string{"unable to create profile"}
	}

	_, err = profile.SaveProfile(db, p)
	if err != nil {
		log.Error().Err(err).Msg("Unable to save new profile")
		return []string{"unable to save profile"}
	}
	log.Debug().Msg("Made it to the end")
	return []string{}
}

// Login takes the email and password, and returns the user ID
func Login(db *sql.DB, email string, password string) (int64, error) {
	log.Debug().
		Str("email", email).
		Str("password", password).
		Msg("Logging in user with")

	if len(email) == 0 {
		log.Debug().
			Int("len(email)", len(email)).
			Msg("Email cannot be blank")
		return 0, errors.New("email cannot be blank")
	}

	if len(password) == 0 {
		log.Debug().
			Int("len(password)", len(email)).
			Msg("Password cannot be blank")
		return 0, errors.New("password cannot be blank")
	}

	q := `	SELECT id, password_hash
			FROM auth_users 
			WHERE email=?`
	res, err := db.Query(q, email)
	if err != nil {
		log.Error().
			Err(err).
			Msg("unable to select from auth_users")
		return 0, err
	}
	defer res.Close()

	var userID int64
	var hash string

	for res.Next() {
		err := res.Scan(&userID, &hash)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("unable to CompareHashAndPassword")
		return 0, err
	}
	log.Debug().
		Int64("userID", userID).
		Msg("user found")
	return userID, nil
}

// Email returns the email for the user
func Email(db *sql.DB, userID int64) (string, error) {
	q := `
	SELECT email
	FROM auth_users
	WHERE id=?`
	res, err := db.Query(q, userID)
	if err != nil {
		return "", err
	}
	defer res.Close()
	var email string
	for res.Next() {
		err := res.Scan(&email)
		if err != nil {
			return "", err
		}
	}
	return email, nil
}

// GetUserByEmail returns the ID of the user by email
func GetUserByEmail(db *sql.DB, email string) (int64, error) {
	q := `
	SELECT id
	FROM auth_users
	WHERE email=?`
	res, err := db.Query(q, email)
	if err != nil {
		return 0, err
	}
	defer res.Close()
	var id int64
	for res.Next() {
		err := res.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

// IsVerified lets us now if the user has verified their email
func IsVerified(db *sql.DB, userID int64) (bool, error) {
	q := `
	SELECT is_verified
	FROM auth_users
	WHERE id=?`
	res, err := db.Query(q, userID)
	if err != nil {
		return false, err
	}
	defer res.Close()
	var isVerified bool
	for res.Next() {
		err := res.Scan(&isVerified)
		if err != nil {
			return false, err
		}
	}
	return isVerified, nil
}

// IsAdmin lets us now if the user is an admin
func IsAdmin(db *sql.DB, userID int64) (bool, error) {
	q := `
	SELECT is_admin
	FROM auth_users
	WHERE id=?`
	res, err := db.Query(q, userID)
	if err != nil {
		return false, err
	}
	defer res.Close()
	var isVerified bool
	for res.Next() {
		err := res.Scan(&isVerified)
		if err != nil {
			return false, err
		}
	}
	return isVerified, nil
}
