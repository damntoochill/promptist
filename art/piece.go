package art

import (
	"database/sql"

	"github.com/rs/zerolog/log"
)

// GetPiece retrieves a piece by ID
func GetPiece(db *sql.DB, id int64) (Piece, bool) {
	q := `
		SELECT piece_id, image_uuid, name, slug, description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces
		WHERE piece_id=?
		`
	res, err := db.Query(q, id)
	if err != nil {
		log.Error().Err(err).Msg("GetPiece: unable to query database")
		return Piece{}, false
	}
	defer res.Close()

	var p Piece

	for res.Next() {
		err := res.Scan(&p.ID, &p.ImageUUID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("GetPiece: unable to scan result")
			return Piece{}, false
		}
	}
	if p.ID == 0 {
		return Piece{}, false
	}
	return p, true
}

// GetPieceByImageUUID retrieves a piece by Image UUID
func GetPieceByImageUUID(db *sql.DB, imageUUIDParam string) (Piece, bool) {
	q := `
		SELECT piece_id, name, slug, description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces
		WHERE image_uuid=?
		`
	res, err := db.Query(q, imageUUIDParam)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return Piece{}, false
	}
	defer res.Close()

	var p Piece
	for res.Next() {
		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return Piece{}, false
		}
	}
	if p.ID == 0 {
		return Piece{}, false
	}
	return p, true
}

// GetPieces from the database
func GetPieces(db *sql.DB) ([]Piece, bool) {
	q := `
		SELECT piece_id, name, slug, description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces
		ORDER BY created_at DESC
		`
	res, err := db.Query(q)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// GetPiecesByUserID from the database
func GetPiecesByUserID(db *sql.DB, userID int64) ([]Piece, bool) {
	q := `
		SELECT piece_id, name, slug, description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces
		WHERE user_id=?
		ORDER BY created_at DESC
		`
	res, err := db.Query(q, userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// GetPiecesByTag from the database
func GetPiecesByTag(db *sql.DB, tagName string) ([]Piece, bool) {
	q := `
		SELECT ap.piece_id, name, slug,  description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces ap
		JOIN art_pieces_tags apt ON apt.piece_id = ap.piece_id
		WHERE tag_name=?
		`
	res, err := db.Query(q, tagName)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// GetPiecesByCollection from the database
func GetPiecesByCollection(db *sql.DB, collectionID int64) ([]Piece, bool) {
	q := `
		SELECT ap.piece_id, name, slug,  description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces ap
		JOIN collection_pieces cp ON cp.piece_id = ap.piece_id
		WHERE collection_id=?
		`
	res, err := db.Query(q, collectionID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// GetPiecesByProgram from the database
func GetPiecesByProgram(db *sql.DB, programID int64) ([]Piece, bool) {
	q := `
		SELECT piece_id, name, slug,  description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces 
		WHERE program_id=?
		`
	res, err := db.Query(q, programID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// GetPiecesByLike from the database
func GetPiecesByLike(db *sql.DB, likerID int64) ([]Piece, bool) {
	q := `
		SELECT ap.piece_id, name, slug,  description, image_uuid, prompt, is_draft, likes, views, saves, comments, al.created_at, program_id, program_name, program_slug, program_cover_image_uuid, ap.user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces ap
		JOIN art_likes al ON al.piece_id = ap.piece_id
		WHERE al.user_id=?
		`
	res, err := db.Query(q, likerID)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// GetPiecesBySearch from the database
func GetPiecesBySearch(db *sql.DB, searchString string) ([]Piece, bool) {
	log.Debug().Str("searchString", searchString).Msg("searching for string")
	q := `
		SELECT piece_id, name, slug,  description, image_uuid, prompt, is_draft, likes, views, saves, comments, created_at, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal
		FROM art_pieces
		WHERE 

		(
			name LIKE ?
			OR description LIKE ?
			OR prompt LIKE ?
			OR tags_literal LIKE ?
		)
		ORDER BY created_at DESC



		`
	res, err := db.Query(q, "%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%")
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return []Piece{}, false
	}
	defer res.Close()

	var pieces []Piece

	for res.Next() {
		var p Piece

		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.ImageUUID, &p.Prompt, &p.IsDraft, &p.Likes, &p.Views, &p.Saves, &p.Comments, &p.CreatedAt, &p.ProgramID, &p.ProgramName, &p.ProgramSlug, &p.ProgramCoverImageUUID, &p.UserID, &p.Username, &p.FullName, &p.ProfilePhotoUUID, &p.TagsLiteral)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return []Piece{}, false
		}
		pieces = append(pieces, p)
	}
	return pieces, true
}

// ValidatePiece creates a new piece without a Piece  ID from the database
func ValidatePiece(userID int64, username string, fullName string, profilePhotoUUID sql.NullString, imageUUID string) (Piece, []string, bool) {
	var problems []string

	if len(imageUUID) < 8 {
		problems = append(problems, "imageUUID must be set")
	}

	if len(problems) > 0 {
		return Piece{}, problems, false
	}

	return Piece{ID: 0, UserID: userID, Username: username, FullName: fullName, ProfilePhotoUUID: profilePhotoUUID, ImageUUID: imageUUID, IsDraft: true}, nil, true
}

// UpdatePiece will save a Piece object and return the piece ID
func UpdatePiece(db *sql.DB, p Piece) (int64, bool) {

	query := `	UPDATE art_pieces 
				SET name=?, slug=?, description=?, image_uuid=?, prompt=?, is_draft=?, likes=?, views=?, saves=?, comments=?, program_id=?, program_name=?, program_slug=?, program_cover_image_uuid=?, user_id=?, username=?, full_name=?, profile_photo_uuid=?, tags_literal=?
				WHERE piece_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query (SavePiece 1)")
		return 0, false
	}
	_, err = stmt.Exec(p.Name, p.Slug, p.Description, p.ImageUUID, p.Prompt, p.IsDraft, p.Likes, p.Views, p.Saves, p.Comments, p.ProgramID, p.ProgramName, p.ProgramSlug, p.ProgramCoverImageUUID, p.UserID, p.Username, p.FullName, p.ProfilePhotoUUID, p.TagsLiteral, p.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	return p.ID, true
}

// AddPiece will save a Piece object and return the piece ID
func AddPiece(db *sql.DB, p Piece) (int64, bool) {

	stmt, err := db.Prepare("INSERT INTO art_pieces (name, slug, description, image_uuid, prompt, is_draft, likes, views, saves, comments, program_id, program_name, program_slug, program_cover_image_uuid, user_id, username, full_name, profile_photo_uuid, tags_literal) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query (SavePiece 2)")
		return 0, false
	}
	res, err := stmt.Exec(p.Name, p.Slug, p.Description, p.ImageUUID, p.Prompt, p.IsDraft, p.Likes, p.Views, p.Saves, p.Comments, p.ProgramID, p.ProgramName, p.ProgramSlug, p.ProgramCoverImageUUID, p.UserID, p.Username, p.FullName, p.ProfilePhotoUUID, p.TagsLiteral)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return 0, false
	}
	pieceID, err := res.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("unable to get piece ID from piece handler")
		return 0, false
	}

	return pieceID, true
}

func UpdateViewCount(db *sql.DB, pieceID int64) bool {
	log.Debug().Msg("UpdateViewCount")
	stmt, err := db.Prepare("UPDATE art_pieces SET views = views + 1 WHERE piece_id = ?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute query")
		return false
	}
	return true
}

// DeletePiece retrieves a piece by ID
func DeletePiece(db *sql.DB, id int64) bool {
	stmt, err := db.Prepare("DELETE FROM art_pieces WHERE piece_id=?")
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare")
		return false
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute")
		return false
	}

	return true
}
