package art

import (
	"database/sql"

	"github.com/rs/zerolog/log"
)

// GetProgramBySlug retrieves a program by ID
func GetProgramBySlug(db *sql.DB, slug string) (Program, bool) {
	q := `
		SELECT program_id, name, slug, description
		FROM art_programs
		WHERE slug=?
		`
	res, err := db.Query(q, slug)
	if err != nil {
		log.Error().Err(err).Msg("unable to query database")
		return Program{}, false
	}
	defer res.Close()

	var p Program

	for res.Next() {
		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description)
		if err != nil {
			log.Error().Err(err).Msg("unable to scan result")
			return Program{}, false
		}
	}
	if p.ID == 0 {
		return Program{}, false
	}
	return p, true
}

// GetProgramByID retrieves a program by ID
func GetProgramByID(db *sql.DB, programID int64) (Program, bool) {
	q := `
		SELECT program_id, name, slug, description
		FROM art_programs
		WHERE program_id=?
		`
	res, err := db.Query(q, programID)
	if err != nil {
		log.Error().Err(err).Msg("GetProgramByID: unable to query database")
		return Program{}, false
	}
	defer res.Close()

	var p Program

	for res.Next() {
		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description)
		if err != nil {
			log.Error().Err(err).Msg("GetProgramByID: unable to scan result")
			return Program{}, false
		}
	}
	if p.ID == 0 {
		return Program{}, false
	}
	return p, true
}

// GetPrograms from the database
func GetPrograms(db *sql.DB) ([]Program, bool) {
	q := `
		SELECT program_id, name, slug, description
		FROM art_programs
		`
	res, err := db.Query(q)
	if err != nil {
		return []Program{}, false
	}
	defer res.Close()

	var programs []Program

	for res.Next() {
		var p Program
		err := res.Scan(&p.ID, &p.Name, &p.Slug, &p.Description)
		if err != nil {
			return []Program{}, false
		}
		programs = append(programs, p)
	}
	return programs, true
}
