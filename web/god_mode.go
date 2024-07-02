package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) godMode(c *gin.Context) {
	c.HTML(http.StatusOK, "god_mode.tmpl", map[string]interface{}{
		"Title": "God Mode",
	})
}
