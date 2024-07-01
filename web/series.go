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

func (s *Server) seriesDetail(c *gin.Context) {
	slug := c.Param("slug")
	series, err := s.db.SeriesSummary(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	items, err := s.db.ItemSummariesBySeries(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "series_detail.tmpl", map[string]interface{}{
		"Title":  series.Name,
		"Series": series,
		"Items":  items,
	})
}
