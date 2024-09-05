package web

import "github.com/gin-gonic/gin"

func (s *Server) flagged(c *gin.Context) {
	scamCreators, err := s.db.ScamCreatorSummaries()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	compromisedSeries, err := s.db.CompromisedSeriesSummaries()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	stolenItems, err := s.db.StolenLostItems()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	s.renderTemplate(c, "flagged.tmpl", map[string]interface{}{
		"Title":             "Flagged",
		"ScamCreators":      scamCreators,
		"CompromisedSeries": compromisedSeries,
		"StolenItems":       stolenItems,
	})
}
