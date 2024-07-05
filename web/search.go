package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) search(c *gin.Context) {
	query := c.Query("q")

	items, err := s.db.Search(query)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "search.tmpl", map[string]interface{}{
		"Title": "Search",
		"Query": query,
		"Items": items,
	})
}
