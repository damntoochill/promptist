package art

import (
	"database/sql"
	"time"

	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// GetLike foo
func GetLike(db *sql.DB, pieceID int64, us *usersession.UserSession) (Like, bool) {
	stmt, err := db.Prepare("SELECT id, piece_id, user_id, created_at FROM art_likes WHERE piece_id = ? AND user_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return Like{}, false
	}
	rows, err := stmt.Query(pieceID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return Like{}, false
	}
	defer rows.Close()
	if !rows.Next() {
		return Like{}, true
	}
	var l Like
	err = rows.Scan(&l.ID, &l.PieceID, &l.UserID, &l.CreatedAt)
	if err != nil {
		log.Error().Err(err).Msg("unable to scan row")
		return Like{}, false
	}
	return l, true
}

// NewLike creates a new follow relationship
func NewLike(db *sql.DB, pieceID int64, us *usersession.UserSession) (Like, bool) {

	_, ok := GetPiece(db, pieceID)
	if !ok {
		return Like{}, false
	}

	insertStmt, err := db.Prepare("INSERT INTO art_likes (piece_id, user_id) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare statement")
		return Like{}, false
	}

	res, err := insertStmt.Exec(pieceID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return Like{}, false
	}

	likeID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get last insert id")
		return Like{}, false
	}

	increaseLikeCount(db, pieceID, us.UserID)

	l := Like{
		ID:        likeID,
		PieceID:   pieceID,
		UserID:    us.UserID,
		CreatedAt: time.Now(),
	}

	return l, true
}

// Unlike creates a new follow relationship
func Unlike(db *sql.DB, pieceID int64, us *usersession.UserSession) bool {

	_, ok := GetPiece(db, pieceID)
	if !ok {
		return false
	}

	stmt, err := db.Prepare("DELETE FROM art_likes WHERE piece_id=? AND user_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare")
		return false
	}

	_, err = stmt.Exec(pieceID, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute")
		return false
	}

	decreaseLikeCount(db, pieceID, us.UserID)

	return true
}

func increaseLikeCount(db *sql.DB, pieceID int64, userID int64) bool {
	log.Debug().Msg("UpdateLikeCount")
	stmt, err := db.Prepare("UPDATE art_pieces SET likes = likes + 1 WHERE piece_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	stmt, err = db.Prepare("UPDATE profile_profiles SET num_likes = num_likes + 1 WHERE user_id = ?")
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

func decreaseLikeCount(db *sql.DB, pieceID int64, userID int64) bool {
	log.Debug().Msg("UpdateLikeCount")
	stmt, err := db.Prepare("UPDATE art_pieces SET likes = likes - 1 WHERE piece_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	stmt, err = db.Prepare("UPDATE profile_profiles SET num_likes = num_likes - 1 WHERE user_id = ?")
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
