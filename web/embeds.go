package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) embedCreator(c *gin.Context) {
	slug := c.Param("slug")
	creator, err := s.db.CreatorSummary(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderNonBaseTemplate(c, "embed_base.tmpl", "embed_creator.tmpl", map[string]interface{}{
		"Title":   creator.Name,
		"Creator": creator,
	})
}
