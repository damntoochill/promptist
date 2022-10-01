package collection

import (
	"database/sql"

	"github.com/promptist/web/art"
	"github.com/promptist/web/usersession"
	"github.com/rs/zerolog/log"
)

// GetCollection returns the an art collection
func GetCollection(db *sql.DB, collectionID int64) (Collection, bool) {
	query := `
		SELECT collection_id, name, description, created_at, is_public, user_id, num_pieces
		FROM collection_collections
		WHERE collection_id=?
		`
	res, err := db.Query(query, collectionID)
	defer res.Close()
	if err != nil {
		log.Error().Err(err).Msg("unable to query follow relationship")
		return Collection{}, false
	}

	var c Collection

	for res.Next() {
		err := res.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.IsPublic, &c.UserID, &c.NumPieces)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return Collection{}, false
		}
	}

	return c, true
}

// GetCollections from the database by a sequence ID
func GetCollections(db *sql.DB, userID int64) ([]Collection, error) {
	q := `
		SELECT collection_id, name, description, created_at, is_public, user_id, num_pieces
		FROM collection_collections
		WHERE user_id=?
		ORDER BY collection_id DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		return []Collection{}, err
	}
	defer res.Close()

	var collections []Collection

	for res.Next() {
		var c Collection
		err := res.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.IsPublic, &c.UserID, &c.NumPieces)
		if err != nil {
			return []Collection{}, err
		}
		collections = append(collections, c)
	}
	return collections, nil
}

// GetCollectionsWithSaved gets the user's collection plus wether or not they have a
// piece saved for a give piece ID
// this is a terrible implementation and needs to be fixed
func GetCollectionsWithSaved(db *sql.DB, userID int64, pieceID int64) ([]Collection, error) {
	q := `
		SELECT collection_id, name, description, created_at, is_public, user_id, num_pieces
		FROM collection_collections
		WHERE user_id=?
		ORDER BY collection_id DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query db")
		return []Collection{}, err
	}
	defer res.Close()

	var collections []Collection

	for res.Next() {
		var c Collection
		err := res.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.IsPublic, &c.UserID, &c.NumPieces)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan db")
			return []Collection{}, err
		}
		p, _ := GetCollectionPiece(db, c.ID, pieceID)
		c.PieceID = p
		collections = append(collections, c)
	}
	return collections, nil
}

// GetCollectionPiece returns the an art collection
func GetCollectionPiece(db *sql.DB, collectionID int64, pieceID int64) (int64, bool) {
	query := `
		SELECT id
		FROM collection_pieces
		WHERE collection_id=? AND piece_id=?
		`
	res, err := db.Query(query, collectionID, pieceID)
	defer res.Close()
	if err != nil {
		log.Error().Err(err).Msg("unable to query follow relationship")
		return 0, false
	}

	var id int64

	for res.Next() {
		err := res.Scan(&id)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan")
			return 0, false
		}
	}

	return id, true
}

// ValidateCollection validates the mandatory fields for a collection
func ValidateCollection(name string) ([]string, bool) {
	var problems []string
	if len(name) < 2 {
		problems = append(problems, "name must be at least 2 characters")
	}

	if len(problems) != 0 {
		return problems, false
	}

	return nil, true
}

// NewCollection creates a new collection
func NewCollection(db *sql.DB, name string, description sql.NullString, isPublic bool, us *usersession.UserSession) (int64, bool) {

	stmt, err := db.Prepare("INSERT INTO collection_collections (name, description, is_public, user_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query (SavePiece 2)")
		return 0, false
	}
	res, err := stmt.Exec(name, description, isPublic, us.UserID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	collectionID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	return collectionID, true
}

// SaveToCollection saves a piece of art to a collection
func SaveToCollection(db *sql.DB, c Collection, p art.Piece, us *usersession.UserSession) (int64, bool) {
	if us.UserID != c.UserID {
		log.Debug().Msg("possbile hacker situation")
		return 0, false
	}

	id, _ := GetCollectionPiece(db, c.ID, p.ID)
	if id != 0 {
		return 0, false
	}

	stmt, err := db.Prepare("INSERT INTO collection_pieces (collection_id, piece_id) VALUES (?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query (SavePiece 2)")
		return 0, false
	}
	res, err := stmt.Exec(c.ID, p.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	insertID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	ok := CountCollectionPieces(db, c)
	if !ok {
		log.Debug().Msg("unable to count pieces")
		return 0, false
	}

	return insertID, true
}

// UpdateCollection updates a profile
func UpdateCollection(db *sql.DB, c Collection) (int64, bool) {

	query := `	UPDATE collection_collections 
				SET name=?, description=?, is_public=?
				WHERE collection_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return 0, false
	}
	_, err = stmt.Exec(c.Name, c.Description, c.IsPublic, c.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	return c.ID, true
}

// CountCollectionPieces updates a profile
func CountCollectionPieces(db *sql.DB, c Collection) bool {

	rows, err := db.Query("SELECT COUNT(*) FROM collection_pieces WHERE collection_id=?", c.ID)
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

	log.Debug().Int64("count", count).Msg("counted this many")

	query := `	UPDATE collection_collections 
				SET num_pieces=?
				WHERE collection_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(count, c.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}
	return true
}

func DeleteCollection(db *sql.DB, collectionID int64, us *usersession.UserSession) bool {

	col, ok := GetCollection(db, collectionID)
	if !ok {
		return false
	}
	if col.UserID != us.UserID {
		return false
	}

	stmt, err := db.Prepare("DELETE FROM collection_pieces WHERE collection_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare")
		return false
	}

	_, err = stmt.Exec(collectionID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute")
		return false
	}

	stmt, err = db.Prepare("DELETE FROM collection_collections WHERE collection_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare")
		return false
	}

	_, err = stmt.Exec(collectionID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute")
		return false
	}

	return true
}

func DeleteCollectionPiece(db *sql.DB, c Collection, pieceID int64, us *usersession.UserSession) bool {

	if c.UserID != us.UserID {
		log.Debug().Msg("possible hacker")
		return false
	}

	stmt, err := db.Prepare("DELETE FROM collection_pieces WHERE collection_id=? AND piece_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare")
		return false
	}

	_, err = stmt.Exec(c.ID, pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute")
		return false
	}

	return true
}
