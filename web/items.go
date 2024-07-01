package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) items(c *gin.Context) {
	items, err := s.db.ItemSummaries()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "items.tmpl", map[string]interface{}{
		"Title": "Items",
		"Items": items,
	})
}

func (s *Server) item(c *gin.Context) {
	sku := c.Param("sku")
	item, err := s.db.ItemSummary(sku)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	itemTransactions, err := s.db.TransactionSummariesByItem(sku)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	addresses, err := s.db.AddressSummariesByItem(sku)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.renderTemplate(c, "item.tmpl", map[string]interface{}{
		"Title":        "Item",
		"Item":         item,
		"Transactions": itemTransactions,
		"Addresses":    addresses,
	})
}
