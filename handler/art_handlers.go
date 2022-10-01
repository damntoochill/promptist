package handler

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/promptist/web/art"
	"github.com/promptist/web/comment"
	"github.com/promptist/web/follow"
	"github.com/promptist/web/image"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/promptist/web/utils"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

// PiecePage is the main jam
func PiecePage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageUUIDParam := c.Param("imageUUID")
		if imageUUIDParam == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		p, ok := art.GetPieceByImageUUID(db, imageUUIDParam)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		var pre []string
		var tags []string
		if p.TagsLiteral.Valid {
			pre = strings.Split(p.TagsLiteral.String, ",")
		}
		for _, tag := range pre {
			tag = slug.Make(tag)
			tag = strings.ToLower(tag)
			tags = append(tags, tag)
		}

		var fo int32
		var myProfile profile.Profile
		var like art.Like
		if us.IsAuthenticated {
			fo = follow.GetFollowOption(db, p.UserID, us.UserID)
			myProfile, _ = profile.GetProfile(db, us.UserID)
			like, _ = art.GetLike(db, p.ID, us)

		}
		comments, _ := comment.GetComments(db, p.ID)
		art.UpdateViewCount(db, p.ID)

		pro, ok := profile.GetProfile(db, p.UserID)

		c.HTML(http.StatusOK, "piece.tmpl", gin.H{
			"UserSession":  us,
			"Piece":        p,
			"Profile":      pro,
			"Tags":         tags,
			"FollowOption": fo,
			"MyProfile":    myProfile,
			"Comments":     comments,
			"Like":         like,
			"Title":        p.Name.String + " " + p.FullName + " on Promptist",
		})
	}
}

// SearchPage is the main jam
func SearchPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Query("q")
		p, ok := art.GetPiecesBySearch(db, query)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"UserSession": us,
			"Pieces":      p,
			"Query":       query,
		})
	}
}

// TagPage helps us navigate to the good stuff
func TagPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		nameParam := c.Param("name")
		log.Debug().Str("name", nameParam).Msg("tag page")
		if nameParam == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		t, ok := art.GetTag(db, nameParam)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		pieces, ok := art.GetPiecesByTag(db, nameParam)
		c.HTML(http.StatusOK, "tag.tmpl", gin.H{
			"UserSession": us,
			"Tag":         t,
			"Pieces":      pieces,
		})
	}
}

// ProgramPage helps us navigate to the good stuff
func ProgramPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		slugParam := c.Param("slug")
		if slugParam == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		p, ok := art.GetProgramBySlug(db, slugParam)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		pieces, ok := art.GetPiecesByProgram(db, p.ID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.HTML(http.StatusOK, "program.tmpl", gin.H{
			"UserSession": us,
			"Program":     p,
			"Pieces":      pieces,
		})
	}
}

// UploadPage for managing your art
func UploadPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.tmpl", gin.H{
			"UserSession": us,
		})
	}
}

// UploadAction uploads a piece of art
func UploadAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := xid.New().String()
		path := "images/" + uuid

		log.Debug().Int64("user_id", us.UserID).Str("username", us.Username).Msg("uploading image")

		if us.IsAuthenticated == false {
			var errors []string
			errors = append(errors, "You must be logged in to upload art")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}

		if us.UserID == 0 {
			var errors []string
			errors = append(errors, "You must be logged in to upload art")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}

		artFile, err := c.FormFile("art")
		if err != nil {
			log.Error().Err(err).Msg("unable to get file from form ")
			var errors []string
			errors = append(errors, "file missing")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
			return
		}

		if artFile == nil {
			c.HTML(http.StatusNotAcceptable, "404.tmpl", gin.H{})
			return
		}

		err = c.SaveUploadedFile(artFile, path)
		if err != nil {
			log.Error().Err(err).Msg("unable to save file to destination ")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		ok := image.Upload(path, path, false, 200, 400, 600, 800, 1000, 1200, 1400, 1600)
		if !ok {
			log.Error().Err(err).Msg("unable to upload file")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		ok = image.Upload(path, path, true, 100, 200, 400, 600)
		if !ok {
			log.Error().Err(err).Msg("unable to upload file")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		err = os.Remove(path)
		if err != nil {
			log.Error().Err(err).Msg("unable to delete file")
		}

		b, ok := profile.GetBrief(db, us.UserID)
		if !ok {
			var errors []string
			errors = append(errors, "unable to get user info")
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": errors,
			})
		}
		log.Debug().Str("username", b.Username).Int64("user id", b.UserID).Msg("profile brief")
		// Create a barebone art piece

		var avatar sql.NullString
		// Set avatar with this piece if it is blank
		if b.PhotoUUID.Valid == false {
			avatar = utils.NewNullString(uuid)
			profile.UpdateProfilePhoto(db, us.UserID, uuid)
			log.Debug().Msg("does not have avatar")

		} else {
			log.Debug().Str("avatar", b.PhotoUUID.String).Msg("had avatar")
			avatar = b.PhotoUUID
		}

		piece := art.Piece{ID: 0, UserID: us.UserID, Username: us.Username, FullName: b.FullName, ProfilePhotoUUID: avatar, ImageUUID: uuid, IsDraft: true}

		_, ok = art.AddPiece(db, piece)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.Redirect(http.StatusFound, "/edit/"+uuid)

	}
}

// EditPage for managing your art
func EditPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageUUIDParam := c.Param("imageUUID")
		if imageUUIDParam == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		p, ok := art.GetPieceByImageUUID(db, imageUUIDParam)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if p.UserID != us.UserID {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		programs, ok := art.GetPrograms(db)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.HTML(http.StatusOK, "edit.tmpl", gin.H{
			"UserSession": us,
			"Piece":       p,
			"Programs":    programs,
		})
	}
}

// EditAction edits a piece of art
func EditAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		idForm := c.PostForm("id")
		nameForm := c.PostForm("name")
		descriptionForm := c.PostForm("description")
		promptForm := c.PostForm("prompt")
		programForm := c.PostForm("program")
		tagsForm := c.PostForm("tags")

		id, err := strconv.ParseInt(idForm, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse id on edit action")
		}
		programID, err := strconv.ParseInt(programForm, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse program id on edit action")
		}
		log.Debug().Int64("id", id).Msg("edit piece ID")

		p, ok := art.GetPiece(db, id)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if p.UserID != us.UserID {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		program, ok := art.GetProgramByID(db, programID)
		if ok {
			log.Debug().Int("program id", int(program.ID)).Msg("we have a program")
			p.ProgramID = program.ID
			p.ProgramName = utils.NewNullString(program.Name)
			p.ProgramSlug = utils.NewNullString(program.Slug)
			p.ProgramCoverImageUUID = utils.NewNullString(program.CoverImageUUID)
		}

		p.Name = utils.NewNullString(nameForm)
		p.Description = utils.NewNullString(descriptionForm)
		p.Prompt = utils.NewNullString(promptForm)
		p.TagsLiteral = utils.NewNullString(tagsForm)

		pieceID, ok := art.UpdatePiece(db, p)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		_ = art.UpdateTags(db, pieceID, tagsForm)

		c.Redirect(http.StatusFound, "/art/"+p.ImageUUID)
	}
}

// DeleteArtAction edits a piece of art
func DeleteArtAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		artIDStr := c.Query("art")

		if artIDStr == "" {
			log.Debug().Msg("no art id")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		artID, err := strconv.ParseInt(artIDStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse id on edit action")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		ok := art.DeletePiece(db, artID)
		if !ok {
			log.Debug().Msg("could not delete")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.Redirect(http.StatusFound, "/account/my-art")
	}
}
