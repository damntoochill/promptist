package follow

import (
	"database/sql"
	"time"

	"github.com/promptist/web/profile"
	"github.com/rs/zerolog/log"
)

// GetFollowOption returns the follower relationship
func GetFollowOption(db *sql.DB, leaderID int64, followerID int64) int32 {
	if leaderID == followerID {
		return 0
	}
	r, ok := GetRelationship(db, leaderID, followerID)
	if !ok {
		return 0
	}
	if r.ID != 0 {
		return 2
	}
	return 1
}

// GetRelationship returns the follower relationship
func GetRelationship(db *sql.DB, leaderID int64, followerID int64) (Relationship, bool) {
	query := `
		SELECT 
			relationship_id, 
			created_at
		FROM follow_relationships
		WHERE leader_id=? AND follower_id=?
		`
	res, err := db.Query(query, leaderID, followerID)
	defer res.Close()
	if err != nil {
		log.Error().Err(err).Msg("unable to query follow relationship")
		return Relationship{}, false
	}

	var r Relationship

	for res.Next() {
		err := res.Scan(&r.ID, &r.CreatedAt)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return Relationship{}, false
		}
	}

	leaderBrief, ok := profile.GetBrief(db, leaderID)
	if !ok {
		return Relationship{}, false
	}
	r.LeaderBrief = leaderBrief

	followerBrief, ok := profile.GetBrief(db, followerID)
	if !ok {
		return Relationship{}, false
	}
	r.FollowerBrief = followerBrief

	return r, true
}

// GetRelationshipByID gets a follower/leader relationship by the relationship id
func GetRelationshipByID(db *sql.DB, relationshipID int64) (Relationship, bool) {
	q := `
		SELECT 
			relationship_id,
			leader_id,
			follower_id, 
			created_at
		FROM follow_relationships
		WHERE relationship_id=?
		`
	res, err := db.Query(q, relationshipID)
	defer res.Close()
	if err != nil {
		log.Error().Err(err).Msg("unable to query follow relationship")
		return Relationship{}, false
	}

	var r Relationship
	var leaderID int64
	var followerID int64

	for res.Next() {
		err := res.Scan(&r.ID, &leaderID, &followerID, &r.CreatedAt)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return Relationship{}, false
		}
	}

	leaderBrief, ok := profile.GetBrief(db, leaderID)
	if !ok {
		return Relationship{}, false
	}
	r.LeaderBrief = leaderBrief

	followerBrief, ok := profile.GetBrief(db, followerID)
	if !ok {
		return Relationship{}, false
	}
	r.FollowerBrief = followerBrief

	return r, true
}

// RelationshipExists checks to see if the follower/leader relationship exists
// return exists bool and ok bool
func RelationshipExists(db *sql.DB, leaderID int64, followerID int64) bool {
	var id int64
	sqlStmt := `SELECT relationship_id FROM follow_relationships WHERE leader_id=? AND follower_id=?`
	err := db.QueryRow(sqlStmt, leaderID, followerID).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Error().Err(err).Msg("unable to  ")

			return false
		}
		log.Error().Err(err).Msg("unable to  ")

		return false
	}

	return true
}

// NewRelationship creates a new follow relationship
func NewRelationship(db *sql.DB, leaderID int64, followerID int64) (Relationship, bool) {

	leaderBrief, ok := profile.GetBrief(db, leaderID)
	if !ok {
		return Relationship{}, false
	}

	followerBrief, ok := profile.GetBrief(db, followerID)
	if !ok {
		return Relationship{}, false
	}

	insertStmt, err := db.Prepare("INSERT INTO follow_relationships (leader_id, follower_id) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare statement")
		return Relationship{}, false
	}

	res, err := insertStmt.Exec(leaderID, followerID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return Relationship{}, false
	}

	followID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get last insert id")
		return Relationship{}, false
	}

	UpdateCounts(db, leaderBrief.UserID)
	UpdateCounts(db, followerBrief.UserID)

	r := Relationship{
		ID:            followID,
		LeaderBrief:   leaderBrief,
		FollowerBrief: followerBrief,
		CreatedAt:     time.Now(),
	}

	return r, true
}

// Unfollow removes the follow/leader relationship from the database
func Unfollow(db *sql.DB, followID int64) bool {
	rel, ok := GetRelationshipByID(db, followID)
	if !ok {
		return false
	}
	stmt, err := db.Prepare("DELETE FROM follow_relationships WHERE relationship_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare")
		return false
	}

	_, err = stmt.Exec(followID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute")
		return false
	}

	UpdateCounts(db, rel.FollowerBrief.UserID)
	UpdateCounts(db, rel.LeaderBrief.UserID)

	return true
}

// GetFollowing from the database
func GetFollowing(db *sql.DB, userID int64) ([]profile.Brief, bool) {
	q := `
		SELECT id, user_id, username, full_name, photo_uuid
		FROM follow_relationships fr
		JOIN profile_profiles pp ON pp.user_id = fr.leader_id
		WHERE follower_id=?
		ORDER BY fr.created_at DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []profile.Brief{}, false
	}
	defer res.Close()

	var briefs []profile.Brief

	for res.Next() {
		var b profile.Brief
		err := res.Scan(&b.ID, &b.UserID, &b.Username, &b.FullName, &b.PhotoUUID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []profile.Brief{}, false
		}
		briefs = append(briefs, b)
	}
	return briefs, true
}

// GetFollowers from the database
func GetFollowers(db *sql.DB, userID int64) ([]profile.Brief, bool) {
	q := `
		SELECT id, user_id, username, full_name, photo_uuid
		FROM follow_relationships fr
		JOIN profile_profiles pp ON pp.user_id = fr.follower_id
		WHERE leader_id=?
		ORDER BY fr.created_at DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []profile.Brief{}, false
	}
	defer res.Close()

	var briefs []profile.Brief

	for res.Next() {
		var b profile.Brief
		err := res.Scan(&b.ID, &b.UserID, &b.Username, &b.FullName, &b.PhotoUUID)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []profile.Brief{}, false
		}
		briefs = append(briefs, b)
	}
	return briefs, true
}

// UpdateCounts updates a profile
func UpdateCounts(db *sql.DB, userID int64) bool {

	rows, err := db.Query("SELECT COUNT(*) FROM follow_relationships WHERE leader_id=?", userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return false
	}
	defer rows.Close()

	var count int64

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return false
		}
	}

	query := `	UPDATE profile_profiles 
				SET num_following=?
				WHERE user_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(count, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	rows, err = db.Query("SELECT COUNT(*) FROM follow_relationships WHERE follower_id=?", userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return false
	}
	defer rows.Close()

	var count2 int64

	for rows.Next() {
		if err := rows.Scan(&count2); err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return false
		}
	}

	query = `	UPDATE profile_profiles 
				SET num_followers=?
				WHERE user_id=?`
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(count, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}

	return true
}
