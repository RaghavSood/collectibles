package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) items(c *gin.Context) {
	items, err := s.db.GetItems()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "items.tmpl", map[string]interface{}{
		"Title": "Items",
		"Items": items,
	})
}
