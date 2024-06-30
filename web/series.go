package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) series(c *gin.Context) {
	series, err := s.db.SeriesSummaries()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "series.tmpl", map[string]interface{}{
		"Title":  "Series",
		"Series": series,
	})
}
