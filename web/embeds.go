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

func (s *Server) embedSeries(c *gin.Context) {
	slug := c.Param("slug")
	series, err := s.db.SeriesSummary(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderNonBaseTemplate(c, "embed_base.tmpl", "embed_series.tmpl", map[string]interface{}{
		"Title":  series.Name,
		"Series": series,
	})
}

func (s *Server) embedItem(c *gin.Context) {
	sku := c.Param("sku")
	item, err := s.db.ItemSummary(sku)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	addresses, err := s.db.AddressSummariesByItem(sku)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderNonBaseTemplate(c, "embed_base.tmpl", "embed_item.tmpl", map[string]interface{}{
		"Title":     "Item",
		"Item":      item,
		"Addresses": addresses,
	})
}
