package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/promptist/web/art"
	"github.com/promptist/web/collection"
	"github.com/promptist/web/profile"
	"github.com/promptist/web/usersession"
	"github.com/promptist/web/utils"
	"github.com/rs/zerolog/log"
)

// CollectionPage is the omega
func CollectionPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameParam := c.Param("username")
		collectionIDParam := c.Param("id")
		collectionID, err := strconv.ParseInt(collectionIDParam, 10, 64)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		p, ok := profile.GetProfileByUsername(db, usernameParam)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		if p.UserID == 0 {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		collection, ok := collection.GetCollection(db, collectionID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		pieces, ok := art.GetPiecesByCollection(db, collectionID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "collection.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
			"Collection":  collection,
			"Pieces":      pieces,
		})
	}
}

// CollectionsPage is the omega
func CollectionsPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		pieceUUIDQuery := c.Query("art")
		var piece art.Piece
		saveMode := false

		usernameParam := c.Param("username")
		p, ok := profile.GetProfileByUsername(db, usernameParam)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}
		if p.UserID == 0 {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please contact customer support"},
			})
			return
		}

		collections, err := collection.GetCollections(db, p.UserID)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		if pieceUUIDQuery != "" {
			piece, ok = art.GetPieceByImageUUID(db, pieceUUIDQuery)
			if ok {
				saveMode = true
			}
		}

		c.HTML(http.StatusOK, "collections.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
			"Collections": collections,
			"Piece":       piece,
			"SaveMode":    saveMode,
		})
	}
}

// CollectionsPromptPage is the omega
func CollectionsPromptPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		pieceUUIDQuery := c.Query("art")
		if pieceUUIDQuery == "" {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		piece, ok := art.GetPieceByImageUUID(db, pieceUUIDQuery)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		collections, err := collection.GetCollectionsWithSaved(db, us.UserID, piece.ID)
		if err != nil {
			log.Debug().Msg("unable to get collections with saved")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "collections_prompt.tmpl", gin.H{
			"UserSession": us,
			"Collections": collections,
			"Piece":       piece,
		})
	}
}

// CollectionNewPage allows the user to edit their profile
func CollectionNewPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		artUUID := c.Query("art")
		saveMode := false

		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please tell customer support"},
			})
			return
		}

		if artUUID != "" {
			saveMode = true
		}

		c.HTML(http.StatusOK, "collection_new.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
			"ArtUUID":     artUUID,
			"SaveMode":    saveMode,
		})
	}
}

// CollectionEditPage allows the user to edit their profile
func CollectionEditPage(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		collectionIDStr := c.Query("collection")
		collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
		col, ok := collection.GetCollection(db, collectionID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if col.ID == 0 {
			log.Debug().Msg("unable to get collection")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if col.UserID != us.UserID {
			log.Debug().Msg("possible hacker action")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		p, ok := profile.GetProfile(db, us.UserID)
		if !ok {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"Errors": []string{"There was a problem loading your info, please tell customer support"},
			})
			return
		}
		c.HTML(http.StatusOK, "collection_edit.tmpl", gin.H{
			"UserSession": us,
			"Profile":     p,
			"Collection":  col,
		})
	}
}

// CollectionEditAction updates the profile
func CollectionEditAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Msg("we are herer")
		nameForm := c.PostForm("name")
		collectionIDStr := c.PostForm("collection-id")
		descriptionForm := c.PostForm("description")
		isPublicForm := c.PostForm("is-public")
		isPublic, err := strconv.ParseBool(isPublicForm)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		description := utils.NewNullString(descriptionForm)
		collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
		col, ok := collection.GetCollection(db, collectionID)
		if !ok {
			log.Debug().Msg("unable to get collection")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		problems, ok := collection.ValidateCollection(nameForm)
		if !ok {
			p, _ := profile.GetProfile(db, us.UserID)
			c.HTML(http.StatusOK, "collection_edit.tmpl", gin.H{
				"UserSession": us,
				"Profile":     p,
				"Problems":    problems,
				"Collection":  col,
			})
		}
		col.Name = nameForm
		col.Description = description
		col.IsPublic = isPublic

		_, ok = collection.UpdateCollection(db, col)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		c.Redirect(http.StatusFound, "/"+us.Username+"/collections/"+collectionIDStr)
	}
}

// CollectionDeleteAction updates the profile
func CollectionDeleteAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectionIDStr := c.Query("collection")
		collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
		ok := collection.DeleteCollection(db, collectionID, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/"+us.Username+"/collections")
	}
}

// CollectionPieceDeleteAction updates the profile
func CollectionPieceDeleteAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectionIDStr := c.Query("collection")
		imageUUID := c.Query("piece")
		collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
		col, ok := collection.GetCollection(db, collectionID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		p, ok := art.GetPieceByImageUUID(db, imageUUID)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		ok = collection.DeleteCollectionPiece(db, col, p.ID, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/"+us.Username+"/collections/"+collectionIDStr)
	}
}

// CollectionNewAction updates the profile
func CollectionNewAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		nameForm := c.PostForm("name")
		descriptionForm := c.PostForm("description")
		isPublicForm := c.PostForm("is-public")
		artUUID := c.PostForm("art-uuid")
		isPublic, err := strconv.ParseBool(isPublicForm)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		description := utils.NewNullString(descriptionForm)

		problems, ok := collection.ValidateCollection(nameForm)
		if !ok {
			p, _ := profile.GetProfile(db, us.UserID)
			c.HTML(http.StatusOK, "collection_new.tmpl", gin.H{
				"UserSession": us,
				"Profile":     p,
				"Problems":    problems,
			})
		}

		collectionID, ok := collection.NewCollection(db, nameForm, description, isPublic, us)
		if !ok {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		if artUUID != "" {
			col, ok := collection.GetCollection(db, collectionID)
			if !ok {
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
			p, ok := art.GetPieceByImageUUID(db, artUUID)
			if !ok {
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
			_, ok = collection.SaveToCollection(db, col, p, us)
			if !ok {
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
				return
			}
		}
		collectionIDStr := strconv.Itoa(int(collectionID))
		c.Redirect(http.StatusFound, "/"+us.Username+"/collections/"+collectionIDStr)
	}
}

// SaveToCollectionAction updates the profile
func SaveToCollectionAction(db *sql.DB, us *usersession.UserSession) gin.HandlerFunc {
	return func(c *gin.Context) {

		collectionIDStr := c.Query("collection")
		imageUUID := c.Query("piece")

		if imageUUID == "" || collectionIDStr == "" {
			log.Debug().Msg("missing query items")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)

		col, ok := collection.GetCollection(db, collectionID)
		if !ok {
			log.Debug().Msg("unable to get collection")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		piece, ok := art.GetPieceByImageUUID(db, imageUUID)
		if !ok {
			log.Debug().Msg("unable to get piece")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}

		_, ok = collection.SaveToCollection(db, col, piece, us)
		if !ok {
			log.Debug().Msg("unable to save collection")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
			return
		}
		c.Redirect(http.StatusFound, "/"+us.Username+"/collections/"+collectionIDStr)
	}
}
