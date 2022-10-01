package royale

import (
	"database/sql"

	"github.com/promptist/web/art"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
	"github.com/sersh88/timeago"
)

// GetRoyales from the database
func GetRoyales(db *sql.DB) ([]Royale, bool) {
	q := `
		SELECT id, name, description, slug, created_at, num_submissions, winner, winning_submission, can_submit, can_vote, prize, round_1_body, round_2_body, round_3_body, round_4_body
		FROM royale_royales
		ORDER BY created_at DESC
		`
	res, err := db.Query(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to query")
		return []Royale{}, false
	}
	defer res.Close()

	var royales []Royale

	for res.Next() {
		var r Royale
		err := res.Scan(&r.ID, &r.Name, &r.Description, &r.Slug, &r.CreatedAt, &r.NumSubmissions, &r.Winner, &r.WinningSubmission, &r.CanSubmit, &r.CanVote, &r.Prize, &r.Round1Body, &r.Round2Body, &r.Round3Body, &r.Round4Body)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return []Royale{}, false
		}
		r.CreatedAtPretty = timeago.New(r.CreatedAt).Format()
		royales = append(royales, r)
	}
	return royales, true
}

// GetRoyaleByID returns foo
func GetRoyaleByID(db *sql.DB, royaleID int64) (Royale, bool) {
	stmt, err := db.Prepare("SELECT id, name, description, slug, created_at, num_submissions, winner, winning_submission, can_submit, can_vote, prize, round_1_body, round_2_body, round_3_body, round_4_body FROM royale_royales WHERE id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Royale{}, false
	}
	rows, err := stmt.Query(royaleID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Royale{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Royale{}, false
	}
	var r Royale
	err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.Slug, &r.CreatedAt, &r.NumSubmissions, &r.Winner, &r.WinningSubmission, &r.CanSubmit, &r.CanVote, &r.Prize, &r.Round1Body, &r.Round2Body, &r.Round3Body, &r.Round4Body)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Royale{}, false
	}
	r.CreatedAtPretty = timeago.New(r.CreatedAt).Format()
	return r, true
}

// GetRoyaleBySlug returns foo
func GetRoyaleBySlug(db *sql.DB, slug string) (Royale, bool) {
	stmt, err := db.Prepare("SELECT id, name, description, slug, created_at, num_submissions, winner, winning_submission, can_submit, can_vote, prize, round_1_body, round_2_body, round_3_body, round_4_body FROM royale_royales WHERE slug = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query GetRoyaleBySlug")
		return Royale{}, false
	}
	rows, err := stmt.Query(slug)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query GetRoyaleBySlug")
		return Royale{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Royale{}, false
	}
	var r Royale
	err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.Slug, &r.CreatedAt, &r.NumSubmissions, &r.Winner, &r.WinningSubmission, &r.CanSubmit, &r.CanVote, &r.Prize, &r.Round1Body, &r.Round2Body, &r.Round3Body, &r.Round4Body)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row GetRoyaleBySlug")
		return Royale{}, false
	}
	r.CreatedAtPretty = timeago.New(r.CreatedAt).Format()
	return r, true
}

// GetSubmissions returns foo
func GetSubmissions(db *sql.DB, royaleID int64) ([]Submission, bool) {
	stmt, err := db.Prepare("SELECT id, royale_id, user_id, art_id, created_at FROM royale_submissions WHERE royale_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query GetSubmissions")
		return nil, false
	}
	res, err := stmt.Query(royaleID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query GetSubmissions")
		return nil, false
	}
	defer res.Close()

	var subs []Submission

	for res.Next() {
		var s Submission
		err := res.Scan(&s.ID, &s.RoyaleID, &s.UserID, &s.ArtID, &s.CreatedAt)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan GetSubmissions")
			return nil, false
		}
		s.CreatedAtPretty = timeago.New(s.CreatedAt).Format()
		subs = append(subs, s)
	}
	return subs, true
}

// Submit returns foo
func Submit(db *sql.DB, royaleID int64, artID string, us *usersession.UserSession) (bool, string) {

	art, ok := art.GetPieceByImageUUID(db, artID)
	if !ok {
		return false, "unable to get art"
	}
	if art.UserID != us.UserID {
		log.Debug().Msg("art user id does not match user session user id, potential hacker")
		return false, "you can't submit someone else's art"
	}

	rows, err := db.Query("SELECT COUNT(*) FROM royale_submissions WHERE royale_id = ? AND user_id = ?", royaleID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query Submit")
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("unable to execute scan Submit")
		}
	}

	if count > 0 {
		return false, "You have already submitted art to this royale."
	}

	stmt, err := db.Prepare("INSERT INTO royale_submissions (royale_id, user_id, art_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false, "there was a server error, plz alert someone"
	}
	_, err = stmt.Exec(royaleID, us.UserID, artID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false, "there was a server error, plz alert someone"
	}
	updateSubmitCount(db, royaleID)

	return true, ""
}

func updateSubmitCount(db *sql.DB, royaleID int64) bool {
	log.Debug().Msg("UpdateSubmitCount")
	stmt, err := db.Prepare("UPDATE royale_royales SET num_submissions = num_submissions + 1 WHERE id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(royaleID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}

// Alert returns foo
func Alert(db *sql.DB, royaleID int64, us *usersession.UserSession) (bool, string) {

	rows, err := db.Query("SELECT COUNT(*) FROM royale_alerts WHERE royale_id = ? AND user_id = ?", royaleID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query Submit")
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("unable to execute scan Submit")
		}
	}

	if count > 0 {
		return false, "You are already judging this royale."
	}

	stmt, err := db.Prepare("INSERT INTO royale_alerts (royale_id, user_id) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false, "there was a server error, plz alert someone"
	}
	_, err = stmt.Exec(royaleID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false, "there was a server error, plz alert someone"
	}

	return true, ""
}

func AmIJudging(db *sql.DB, royaleID int64, us *usersession.UserSession) bool {
	rows, err := db.Query("SELECT COUNT(*) FROM royale_alerts WHERE royale_id = ? AND user_id = ?", royaleID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query Submit")
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("unable to execute scan Submit")
		}
	}

	if count > 0 {
		return true
	}

	return false
}
