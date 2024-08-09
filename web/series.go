package web

import (
	"fmt"
	"net/http"

	"github.com/RaghavSood/collectibles/notes"
	"github.com/RaghavSood/collectibles/types"
	"github.com/RaghavSood/collectibles/util"
	"github.com/gin-gonic/gin"
)

func (s *Server) series(c *gin.Context) {
	series, err := s.db.SeriesSummaries()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for i := range series {
		creators, err := s.db.GetCreatorsBySeries(series[i].Slug)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		series[i].Creators = creators
	}

	s.renderTemplate(c, "series.tmpl", map[string]interface{}{
		"Title":  "Series",
		"Series": series,
	})
}

func (s *Server) seriesDetail(c *gin.Context) {
	slug := c.Param("slug")
	series, err := s.db.SeriesSummary(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	creators, err := s.db.GetCreatorsBySeries(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	series.Creators = creators

	itemSummaries, err := s.db.ItemAddressSummariesBySeries(slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	flags, err := s.db.GetFlags(types.FLAG_SCOPE_SERIES, series.Slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notePointer := notes.NotePointer{
		NoteType:     notes.Series,
		PathElements: []string{series.Slug},
	}

	notes := notes.ReadNotes([]notes.NotePointer{notePointer})

	s.renderTemplate(c, "series_detail.tmpl", map[string]interface{}{
		"Title":         series.Name,
		"Desc":          fmt.Sprintf("%s has %d items worth %s BTC (%s USD)", series.Name, series.ItemCount, series.TotalValue.SatoshisToBTC(true), util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(series.TotalValue)))),
		"Series":        series,
		"Notes":         notes,
		"ItemSummaries": itemSummaries,
		"Flags":         flags,
	})
}
