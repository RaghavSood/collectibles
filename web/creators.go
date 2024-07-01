package web

import (
	"database/sql"
	"net/http"

	"github.com/RaghavSood/collectibles/notes"
	"github.com/gin-gonic/gin"
)

func (s *Server) creators(c *gin.Context) {
	creators, err := s.db.CreatorSummaries()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "creators.tmpl", map[string]interface{}{
		"Title":    "Creators",
		"Creators": creators,
	})
}

func (s *Server) creator(c *gin.Context) {
	slug := c.Param("slug")
	creator, err := s.db.CreatorSummary(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	series, err := s.db.SeriesSummariesByCreator(slug)
	if err != nil && err != sql.ErrNoRows {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notePointer := notes.NotePointer{
		NoteType:     notes.Creator,
		PathElements: []string{creator.Slug},
	}

	notes := notes.ReadNotes([]notes.NotePointer{notePointer})

	s.renderTemplate(c, "creator.tmpl", map[string]interface{}{
		"Title":   creator.Name,
		"Creator": creator,
		"Series":  series,
		"Notes":   notes,
	})
}
