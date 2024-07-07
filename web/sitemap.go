package web

import (
	"net/http"
	"strconv"
	"strings"

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

func (s *Server) sitemapItems(c *gin.Context) {
	index := c.Param("index")

	parts := strings.Split(index, ".")
	if len(parts) != 2 || parts[1] != "xml" {
		c.AbortWithStatus(404)
		return
	}

	page, err := strconv.Atoi(parts[0])
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	if page < 0 {
		c.AbortWithStatus(400)
		return
	}

	pageSize := 25000
	offset := page * pageSize
	items, err := s.db.GetItemPage(pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	si := sitemap.New()
	for _, item := range items {
		result := "https://collectible.money/item/" + item.SKU
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(http.StatusOK, si)
}

func (s *Server) sitemapItemIndex(c *gin.Context) {
	stats, err := s.db.GeneralStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch general statistics"})
		return
	}
	itemCount := stats.Items
	pageSize := 25000

	pages := (itemCount / pageSize) + 1
	si := sitemap.NewSitemapIndex()
	for i := 0; i < pages; i++ {
		result := "https://collectible.money/sitemap/items/" + strconv.Itoa(i) + ".xml"
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(http.StatusOK, si)
}
