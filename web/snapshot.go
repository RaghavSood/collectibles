package web

import (
	"context"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

func (s *Server) snapshot(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	path = path[1:] // Remove the leading slash

	timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := chromedp.NewContext(
		timeoutCtx,
	)
	// to release the browser resources when
	// it is no longer needed
	defer cancel()

	var screenshotBuffer []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:8086/embed/"+path),
		chromedp.Screenshot("#embed", &screenshotBuffer, chromedp.NodeVisible),
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "image/png", screenshotBuffer)
}
