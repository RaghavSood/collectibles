package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snabb/sitemap"
)

func (s *Server) sitemapCreators(c *gin.Context) {
	creators, err := s.db.GetCreators()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch creators"})
		return
	}

	si := sitemap.New()
	for _, creator := range creators {
		result := "https://collectible.money/creator/" + creator.Slug
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(http.StatusOK, si)
}

func (s *Server) sitemapSeries(c *gin.Context) {
	series, err := s.db.GetSeries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch series"})
		return
	}

	si := sitemap.New()
	for _, serie := range series {
		result := "https://collectible.money/series/" + serie.Slug
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(http.StatusOK, si)
}
