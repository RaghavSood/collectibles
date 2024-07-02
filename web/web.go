package web

import (
	"net/http"
	"os"

	"github.com/RaghavSood/collectibles/clogger"
	"github.com/RaghavSood/collectibles/middleware"
	"github.com/RaghavSood/collectibles/static"
	"github.com/RaghavSood/collectibles/storage"
	"github.com/RaghavSood/collectibles/templates"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db       storage.Storage
	readOnly bool
}

func NewServer(db storage.Storage, noindex bool) *Server {
	return &Server{
		db:       db,
		readOnly: noindex,
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

	s.renderTemplate(c, "index.tmpl", map[string]interface{}{
		"Title": "Home",
		"Stats": stats,
	})
}

func (s *Server) about(c *gin.Context) {
	s.renderTemplate(c, "about.tmpl", map[string]interface{}{
		"Title": "About",
	})
}

func (s *Server) renderTemplate(c *gin.Context, template string, params map[string]interface{}) {
	tmpl := templates.New()
	err := tmpl.Render(c.Writer, template, params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
