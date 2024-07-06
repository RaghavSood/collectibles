package web

import (
	"fmt"
	"time"

	"github.com/RaghavSood/collectibles/types"
	"github.com/RaghavSood/collectibles/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func (s *Server) feedCreator(c *gin.Context) {
	creator := c.Param("creator")
	transactions, err := s.db.TransactionSummariesByCreator(creator, 100)
	if err != nil {
		c.XML(500, gin.H{"error": "Error fetching transactions"})
		return
	}

	feed := &feeds.Feed{
		Title:       fmt.Sprintf("Collectible - Creator %s", creator),
		Link:        &feeds.Link{Href: "https://collectible.money/creator/" + creator},
		Description: fmt.Sprintf("Transactions involving %s", creator),
		Author:      &feeds.Author{Name: "Collectible", Email: "hello@collectible.money"},
		Created:     time.Now(),
	}

	feed.Items = feedItems(transactions)

	rss, err := feed.ToAtom()
	if err != nil {
		c.XML(500, gin.H{"error": "Error generating feed"})
		return
	}

	c.Data(200, "application/atom+xml", []byte(rss))
}

func (s *Server) feedSeries(c *gin.Context) {
	series := c.Param("series")
	transactions, err := s.db.TransactionSummariesBySeries(series, 100)
	if err != nil {
		c.XML(500, gin.H{"error": "Error fetching transactions"})
		return
	}

	feed := &feeds.Feed{
		Title:       fmt.Sprintf("Collectible - Series %s", series),
		Link:        &feeds.Link{Href: "https://collectible.money/series/" + series},
		Description: fmt.Sprintf("Transactions involving %s", series),
		Author:      &feeds.Author{Name: "Collectible", Email: "hello@collectible.money"},
		Created:     time.Now(),
	}

	feed.Items = feedItems(transactions)

	rss, err := feed.ToAtom()
	if err != nil {
		c.XML(500, gin.H{"error": "Error generating feed"})
		return
	}

	c.Data(200, "application/atom+xml", []byte(rss))
}

func (s *Server) feedItem(c *gin.Context) {
	sku := c.Param("item")
	transactions, err := s.db.TransactionSummariesByItem(sku)
	if err != nil {
		c.XML(500, gin.H{"error": "Error fetching transactions"})
		return
	}

	item, err := s.db.ItemSummary(sku)
	if err != nil {
		c.XML(500, gin.H{"error": "Error fetching item"})
		return
	}

	title := fmt.Sprintf("Collectible - %s", item.SeriesName)
	if item.Serial != "" {
		title = fmt.Sprintf("Collectible - %s (%s)", item.SeriesName, item.Serial)
	}

	description := fmt.Sprintf("Transactions involving %s", item.SeriesName)
	if item.Serial != "" {
		description = fmt.Sprintf("Transactions involving %s (%s)", item.SeriesName, item.Serial)
	}

	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: "https://collectible.money/item/" + sku},
		Description: description,
		Author:      &feeds.Author{Name: "Collectible", Email: "hello@collectible.money"},
		Created:     time.Now(),
	}

	feed.Items = feedItems(transactions)

	rss, err := feed.ToAtom()
	if err != nil {
		c.XML(500, gin.H{"error": "Error generating feed"})
		return
	}

	c.Data(200, "application/atom+xml", []byte(rss))
}

func (s *Server) feedAll(c *gin.Context) {
	recentTransactions, err := s.db.TransactionSummaries(100)
	if err != nil {
		c.XML(500, gin.H{"error": "Error fetching recent transactions"})
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Collectible - Recent Transactions",
		Link:        &feeds.Link{Href: "https://collectible.money/"},
		Description: "Recent transactions involving Bitcoin and other crypto collectibles",
		Author:      &feeds.Author{Name: "Collectible", Email: "hello@collectible.money"},
		Created:     now,
	}

	feed.Items = feedItems(recentTransactions)

	rss, err := feed.ToAtom()
	if err != nil {
		c.XML(500, gin.H{"error": "Error generating feed"})
		return
	}

	c.Data(200, "application/atom+xml", []byte(rss))
}

func feedItems(txs []types.Transaction) []*feeds.Item {
	var items []*feeds.Item

	for _, tx := range txs {
		title := fmt.Sprintf("An %s transaction for %s", tx.TransactionType, tx.SeriesName)
		if tx.Serial != "" {
			title = fmt.Sprintf("An %s transaction for %s (%s)", tx.TransactionType, tx.SeriesName, tx.Serial)
		}

		var description string
		usdValue := fmt.Sprintf("%.2f", util.BTCValueToUSD(tx.Value))
		if tx.TransactionType == "incoming" {
			description = fmt.Sprintf("Received %s BTC (%s USD)", tx.Value.SatoshisToBTC(true), util.FormatNumber(usdValue))
		} else {
			description = fmt.Sprintf("Sent %s BTC (%s USD)", tx.Value.SatoshisToBTC(true), util.FormatNumber(usdValue))
		}
		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://collectible.money/item/%s", tx.SKU)},
			Description: description,
			Author:      &feeds.Author{Name: "Collectible", Email: "hello@collectible.money"},
			Created:     tx.BlockTime.UTC(),
		}

		items = append(items, item)
	}

	return items
}
