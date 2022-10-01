package profile

import (
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
)

// GetProfile retrieves a profile from the database by a UserID
func GetProfile(db *sql.DB, userID int64) (Profile, bool) {
	q := `
		SELECT id, user_id, username, full_name, bio, location, photo_uuid, num_following, num_followers, num_collections, num_pieces, num_likes, num_views, num_comments, is_new
		FROM profile_profiles
		WHERE user_id=?
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return Profile{}, false
	}
	defer res.Close()

	var p Profile

	for res.Next() {
		err := res.Scan(&p.ID, &p.UserID, &p.Username, &p.FullName, &p.Bio, &p.Location, &p.PhotoUUID, &p.NumFollowing, &p.NumFollowers, &p.NumCollections, &p.NumPieces, &p.NumLikes, &p.NumViews, &p.NumComments, &p.IsNew)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return Profile{}, false
		}
	}
	return p, true
}

// GetProfileByUsername retrieves a profile from the database by a Username
func GetProfileByUsername(db *sql.DB, username string) (Profile, bool) {
	q := `
		SELECT id, user_id, username, full_name, bio, location, photo_uuid, num_following, num_followers, num_collections, num_pieces, num_likes, num_views, num_comments, is_new
		FROM profile_profiles
		WHERE username=?
		`
	res, err := db.Query(q, username)
	if err != nil {
		log.Error().Err(err).Msg("unable to query db")
		return Profile{}, false
	}
	defer res.Close()

	var p Profile

	for res.Next() {
		err := res.Scan(&p.ID, &p.UserID, &p.Username, &p.FullName, &p.Bio, &p.Location, &p.PhotoUUID, &p.NumFollowing, &p.NumFollowers, &p.NumCollections, &p.NumPieces, &p.NumLikes, &p.NumViews, &p.NumComments, &p.IsNew)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan db")
			return Profile{}, false
		}
	}
	return p, true
}

// GetBrief retrieves a profile from the database by a UserID
func GetBrief(db *sql.DB, userID int64) (Brief, bool) {
	log.Debug().Int64("userID", userID).Msg("get brief by ID")
	q := `
		SELECT id, user_id, username, full_name, photo_uuid
		FROM profile_profiles
		WHERE user_id=?
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return Brief{}, false
	}
	defer res.Close()

	var p Brief

	for res.Next() {
		err := res.Scan(&p.ID, &p.UserID, &p.Username, &p.FullName, &p.PhotoUUID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan row")
			return Brief{}, false
		}
	}
	return p, true
}

// GetProfiles from the database by a sequence ID
func GetProfiles(db *sql.DB) ([]Profile, bool) {
	q := `
		SELECT id, user_id, username, full_name, bio, location, photo_uuid, num_following, num_followers, num_collections, num_pieces, num_likes, num_views, num_comments
		FROM profile_profiles
		ORDER BY id DESC
		`
	res, err := db.Query(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Profile{}, false
	}
	defer res.Close()

	var profiles []Profile

	for res.Next() {
		var p Profile
		err := res.Scan(&p.ID, &p.UserID, &p.Username, &p.FullName, &p.Bio, &p.Location, &p.PhotoUUID, &p.NumFollowing, &p.NumFollowers, &p.NumCollections, &p.NumPieces, &p.NumLikes, &p.NumViews, &p.NumComments)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan ")
			return []Profile{}, false
		}
		profiles = append(profiles, p)
	}
	return profiles, true
}

// GetPioneers from the database by a sequence ID
func GetPioneers(db *sql.DB) ([]Profile, bool) {
	q := `
		SELECT id, user_id, username, full_name, bio, location, photo_uuid, num_following, num_followers, num_collections, num_pieces, num_likes, num_views, num_comments
		FROM profile_profiles
		WHERE photo_uuid IS NOT NULL
		ORDER BY id
		LIMIT 1000
		`
	res, err := db.Query(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Profile{}, false
	}
	defer res.Close()

	var profiles []Profile

	for res.Next() {
		var p Profile
		err := res.Scan(&p.ID, &p.UserID, &p.Username, &p.FullName, &p.Bio, &p.Location, &p.PhotoUUID, &p.NumFollowing, &p.NumFollowers, &p.NumCollections, &p.NumPieces, &p.NumLikes, &p.NumViews, &p.NumComments)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan ")
			return []Profile{}, false
		}
		profiles = append(profiles, p)
	}
	return profiles, true
}

// ValidateProfile validates the mandatory fields for a profile
func ValidateProfile(username string, fullName string) ([]string, bool) {
	var problems []string
	if len(username) < 3 {
		problems = append(problems, "username must be at least 3 characters")
	}

	if len(fullName) < 3 {
		problems = append(problems, "full name must be at least 3 characters")
	}

	if len(problems) != 0 {
		return problems, false
	}

	return nil, true
}

// NewProfile creates a new profile without a Profile  ID from the database
func NewProfile(userID int64, username string, fullName string, bio sql.NullString, location sql.NullString) (Profile, error) {
	if len(username) < 3 {
		err := errors.New("username must be at least 3 characters")
		log.Error().Err(err).Msg("username must be at least 3 characters")
		return Profile{}, err
	}

	if len(fullName) < 3 {
		err := errors.New("full name must be at least 3 characters")
		log.Error().Err(err).Msg("full name must be at least 3 characters")
		return Profile{}, err
	}

	return Profile{ID: 0, UserID: userID, Username: username, FullName: fullName, Bio: bio, Location: location}, nil
}

// SaveProfile will save a Profile object and return the profile ID
func SaveProfile(db *sql.DB, p Profile) (int64, error) {

	// Query the database to see if a profile already exists for UserID
	var rowExists bool
	query := `SELECT id FROM profile_profiles WHERE user_id=?`
	err := db.QueryRow(query, p.UserID).Scan(&rowExists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal().Err(err).Msg("error checking if row exists '%s' %v")
		return 0, err
	}
	if rowExists {
		// Update current profile
		query := `	UPDATE profile_profiles 
				  	SET username=?, full_name=?, bio=?, location=?
				  	WHERE user_id=?`
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare query")
			return 0, err
		}
		_, err = stmt.Exec(p.Username, p.FullName, p.Bio, p.Location, p.UserID)
		if err != nil {
			log.Error().Err(err).Msg("unable to execute statement")
			return 0, err
		}
		return p.ID, nil
	}

	// Insert new profile
	stmt, err := db.Prepare("INSERT INTO profile_profiles (user_id, username, full_name, bio, location) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, err
	}
	res, err := stmt.Exec(p.UserID, p.Username, p.FullName, p.Bio, p.Location)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, err
	}
	profileID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get post ID from post handler")
		return 0, err
	}
	return profileID, nil
}

// UpdateProfile updates a profile
func UpdateProfile(db *sql.DB, p Profile) (int64, error) {

	query := `	UPDATE profile_profiles 
				SET username=?, full_name=?, bio=?, location=?
				WHERE id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, err
	}
	_, err = stmt.Exec(p.Username, p.FullName, p.Bio, p.Location, p.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, err
	}
	ok := MassUpdate(db, p.UserID)
	if !ok {
		log.Error().Err(err).Msg("unable to update user")
		return 0, err
	}
	return p.ID, nil
}

// UpdateProfilePhoto updates a profile
func UpdateProfilePhoto(db *sql.DB, userID int64, photoUUID string) bool {

	query := `	UPDATE profile_profiles 
				SET photo_uuid=?
				WHERE user_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query (UpdateProfilePhoto)")
		return false
	}
	_, err = stmt.Exec(photoUUID, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement (UpdateProfilePhoto)")
		return false
	}
	ok := MassUpdate(db, userID)
	if !ok {
		log.Error().Err(err).Msg("unable to update profile (UpdateProfilePhoto)")
		return false
	}
	return true
}

func MassUpdate(db *sql.DB, userID int64) bool {
	profile, ok := GetProfile(db, userID)
	if !ok {
		return false
	}
	query := `	UPDATE art_pieces 
				SET full_name=?, username=?, profile_photo_uuid=?
				WHERE user_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(profile.FullName, profile.Username, profile.PhotoUUID, profile.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	query = `	UPDATE comment_comments 
				SET full_name=?, username=?, avatar=?
				WHERE user_id=?`
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(profile.FullName, profile.Username, profile.PhotoUUID, profile.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	query = `	UPDATE comment_pro_comments 
	SET full_name=?, username=?, avatar=?
	WHERE user_id=?`
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(profile.FullName, profile.Username, profile.PhotoUUID, profile.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	query = `	UPDATE notif_notifications 
	SET from_full_name=?, from_username=?, from_avatar=?
	WHERE from_user_id=?`
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(profile.FullName, profile.Username, profile.PhotoUUID, profile.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	query = `	UPDATE post_posts 
	SET full_name=?, username=?, avatar=?
	WHERE user_id=?`
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(profile.FullName, profile.Username, profile.PhotoUUID, profile.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	return true
}
