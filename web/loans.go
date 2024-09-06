package web

import (
	"github.com/RaghavSood/collectibles/prices"
	"github.com/gin-gonic/gin"
)

func (s *Server) loans(c *gin.Context) {
	price, err := prices.GetBTCUSDPrice()
	if err != nil {
		price = 0
	}

	s.renderTemplate(c, "loans.tmpl", map[string]interface{}{
		"Title": "Loans",
		"Price": price,
	})
}
