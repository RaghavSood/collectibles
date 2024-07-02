package web

import (
	"github.com/RaghavSood/collectibles/templates"
	"github.com/gin-gonic/gin"
)

func (s *Server) godMode(c *gin.Context) {
	templates.RenderSingle(c.Writer, "god_mode.tmpl", nil)
}
