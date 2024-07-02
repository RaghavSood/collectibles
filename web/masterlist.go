package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) masterList(c *gin.Context) {
	masterlist, err := s.db.GodView()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "masterlist.tmpl", map[string]interface{}{
		"Title":      "Masterlist",
		"Masterlist": masterlist,
	})
}
