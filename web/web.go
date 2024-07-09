package web

import (
	"net/http"
	"os"
	"time"

	"github.com/RaghavSood/collectibles/clogger"
	"github.com/RaghavSood/collectibles/middleware"
	"github.com/RaghavSood/collectibles/static"
	"github.com/RaghavSood/collectibles/storage"
	"github.com/RaghavSood/collectibles/templates"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/rs/zerolog/log"
)

type Server struct {
	db            storage.Storage
	readOnly      bool
	snapshotCache *expirable.LRU[string, []byte]
}

func NewServer(db storage.Storage, noindex bool) *Server {
	cache := expirable.NewLRU[string, []byte](500, nil, 90*time.Second)
	return &Server{
		db:            db,
		readOnly:      noindex,
		snapshotCache: cache,
	}
}

func (s *Server) Serve() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.StructuredLogger(clogger.NewLogger("gin-webui")))
	router.Use(gin.Recovery())

	router.GET("/", s.index)
	router.GET("/about", s.about)

	router.GET("/creators", s.creators)
	router.GET("/creator/:slug", s.creator)

	router.GET("/series", s.series)
	router.GET("/series/:slug", s.seriesDetail)

	router.GET("/items", s.items)
	router.GET("/item/:sku", s.item)

	router.GET("/ogimage/:slug", s.ogimage)

	router.GET("/godmode", s.godMode)

	router.GET("/masterlist", s.masterList)

	router.GET("/search", s.search)

	router.GET("/snapshot/*path", s.snapshot)
	embeds := router.Group("/embed")
	{
		embeds.GET("/creator/:slug", s.embedCreator)
		embeds.GET("/series/:slug", s.embedSeries)
		embeds.GET("/item/:sku", s.embedItem)
	}

	feeds := router.Group("/feeds")
	{
		feeds.GET("/creator/:creator", s.feedCreator)
		feeds.GET("/series/:series", s.feedSeries)
		feeds.GET("/item/:item", s.feedItem)
		feeds.GET("/all", s.feedAll)
	}

	router.StaticFS("/static", http.FS(static.Static))
	// Serve /favicon.ico and /robots.txt from the root
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(static.Static))
	})

	router.GET("/robots.txt", func(c *gin.Context) {
		c.FileFromFS("robots.txt", http.FS(static.Static))
	})

	sitemap := router.Group("/sitemap")
	{
		sitemap.GET("/creators.xml", s.sitemapCreators)
		sitemap.GET("/series.xml", s.sitemapSeries)

		sitemap.GET("/index/items.xml", s.sitemapItemIndex)
		sitemap.GET("/items/:index", s.sitemapItems)
	}

	port := os.Getenv("COLLECTIBLES_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

func (s *Server) index(c *gin.Context) {
	stats, err := s.db.GeneralStatistics()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	recentRedemptions, err := s.db.RecentRedemptions(50)

	s.renderTemplate(c, "index.tmpl", map[string]interface{}{
		"Title":             "Home",
		"Stats":             stats,
		"RecentRedemptions": recentRedemptions,
	})
}

func (s *Server) about(c *gin.Context) {
	s.renderTemplate(c, "about.tmpl", map[string]interface{}{
		"Title": "About",
	})
}

func (s *Server) renderNonBaseTemplate(c *gin.Context, newBase, template string, params map[string]interface{}) {
	tmpl := templates.New()
	err := tmpl.RenderNonBase(c.Writer, newBase, template, params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) renderTemplate(c *gin.Context, template string, params map[string]interface{}) {
	importOngoing := false
	sQueueLen, txQueueLen, err := s.db.GetQueueStats()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get queue stats")
	}

	if sQueueLen > 0 || txQueueLen > 0 {
		importOngoing = true
	}

	params["ImportOngoing"] = importOngoing
	params["ScriptQueueLen"] = sQueueLen
	params["TxQueueLen"] = txQueueLen

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, template, params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
