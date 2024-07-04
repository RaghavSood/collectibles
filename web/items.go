package web

import (
	"fmt"
	"net/http"

	"github.com/RaghavSood/collectibles/util"
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
		"Title":        fmt.Sprintf("%s - %s", item.SeriesName, item.SerialString()),
		"Desc":         fmt.Sprintf("An item in the %s series holding %s BTC (%s USD)", item.SeriesName, item.TotalValue.SatoshisToBTC(true), util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(item.TotalValue)))),
		"Item":         item,
		"Transactions": itemTransactions,
		"Addresses":    addresses,
	})
}
