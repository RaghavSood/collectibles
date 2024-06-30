package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) creators(c *gin.Context) {
	creators, err := s.db.GetCreators()
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
	creator, err := s.db.GetCreator(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "creator.tmpl", map[string]interface{}{
		"Title":   creator.Name,
		"Creator": creator,
	})
}
