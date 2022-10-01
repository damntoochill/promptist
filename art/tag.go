package art

import (
	"database/sql"
	"strings"

	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
)

// GetTag retrieves a product from the database by name (slugified)
func GetTag(db *sql.DB, name string) (Tag, bool) {
	q := `
		SELECT tag_id, name, description, total
		FROM art_tags
		WHERE name=?
		`
	res, err := db.Query(q, name)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return Tag{}, false
	}
	defer res.Close()

	var t Tag

	for res.Next() {
		err := res.Scan(&t.ID, &t.Name, &t.Description, &t.Total)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return Tag{}, false
		}
	}
	if t.ID == 0 {
		log.Debug().Msg("The id of the tag is 0")
		return Tag{}, false
	}
	return t, true
}

// UpdateTags updates the tags for a piece
func UpdateTags(db *sql.DB, pieceID int64, tagsLiteral string) bool {

	_, err := db.Exec("DELETE FROM art_pieces_tags WHERE piece_id = ?", pieceID)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete rows from art_pieces_tags")
		return false
	}

	if tagsLiteral != "" {
		tags := ProcessTags(tagsLiteral)

		for _, tag := range tags {

			t, _ := GetTag(db, tag)
			if t.ID == 0 {
				// insert
				stmt, err := db.Prepare("INSERT INTO art_tags (name) VALUES (?)")
				if err != nil {
					log.Error().Err(err).Msg("unable to prepare query")
					return false
				}
				_, err = stmt.Exec(tag)
				if err != nil {
					log.Error().Err(err).Msg("unable to execute statement")
					return false
				}
			} else {
				// update count
				t.Total = t.Total + 1
				UpdateTag(db, t)
			}

			stmt, err := db.Prepare("INSERT INTO art_pieces_tags (tag_name, piece_id) VALUES (?, ?)")
			if err != nil {
				log.Error().Err(err).Msg("unable to prepare query")
				return false
			}
			_, err = stmt.Exec(tag, pieceID)
			if err != nil {
				log.Error().Err(err).Msg("unable to execute statement")
				return false
			}
		}
	}
	return true
}

// ProcessTags will do some things to tags
func ProcessTags(tagsLiteral string) []string {
	var pre []string
	var tags []string
	pre = strings.Split(tagsLiteral, ",")
	for _, tag := range pre {
		tag = slug.Make(tag)
		tag = strings.ToLower(tag)
		tags = append(tags, tag)
	}
	return tags
}

// UpdateTag updates a profile
func UpdateTag(db *sql.DB, t Tag) bool {

	query := `	UPDATE art_tags 
				SET name=?, description=?, total=?
				WHERE tag_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare query")
		return false
	}
	_, err = stmt.Exec(t.Name, t.Description, t.Total, t.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to execute statement")
		return false
	}
	return true
}
